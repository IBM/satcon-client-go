package actions_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"text/template"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.ibm.com/coligo/satcon-client/client/actions"
)

var _ = Describe("Query", func() {
	var (
		argMap map[string]string
	)

	BeforeEach(func() {
		argMap = map[string]string{
			"orgId":     "String!",
			"flavor":    "String!",
			"dimension": "JSON!",
		}
	})

	Describe("BuildArgsList", func() {
		It("Returns a string containing a list delimited by ', '", func() {
			argList := BuildArgsList(argMap)
			// We cannot just do a simple string comparision, because go does not automatically
			// sort map keys, nor do we really want to require them to be sorted.
			// So we tokenize the returned string and make sure it has the right elements and
			// has no trailing comma/whitespace.
			tokens := strings.Split(argList, ", ")
			Expect(tokens).To(ConsistOf(
				"$orgId: String!",
				"$flavor: String!",
				"$dimension: JSON!",
			))

			Expect(argList).NotTo(HaveSuffix(", "))
		})
	})

	Describe("BuildArgVarsList", func() {
		It("Returns a correct GraphQL string for the argument variables", func() {
			argVarList := BuildArgVarsList(argMap)

			tokens := strings.Split(argVarList, ", ")
			Expect(tokens).To(ConsistOf(
				"orgId: $orgId",
				"flavor: $flavor",
				"dimension: $dimension",
			))

			Expect(argVarList).NotTo(HaveSuffix(", "))
		})
	})

	Describe("BuildRequest", func() {
		var (
			endpoint, token string
			payload         *bytes.Buffer
		)

		BeforeEach(func() {
			endpoint = "http://foo.bar"
			token = "atoken"
			payload = bytes.NewBuffer([]byte("stringifiedbody"))
		})

		It("Returns a valid request instance", func() {
			req := BuildRequest(payload, endpoint, token)
			Expect(req).NotTo(BeNil())
		})

		It("Populates the request with the required headers", func() {
			req := BuildRequest(payload, endpoint, token)
			Expect(req.Header).To(HaveKeyWithValue(
				MatchRegexp(`[cC]ontent-[tT]ype`),
				ContainElement("application/json"),
			))
			Expect(req.Header).To(HaveKeyWithValue(
				MatchRegexp("[aA]uthorization"),
				ContainElement(fmt.Sprintf("Bearer %s", token)),
			))
		})
	})

	Describe("BuildRequestBody", func() {
		type requestVars struct {
			GraphQLQuery
			First string
			Last  string
		}

		var (
			requestTemplate string
			vars            requestVars
			funcs           template.FuncMap
		)

		BeforeEach(func() {
			requestTemplate = `{{define "vars"}}"first":"{{js .First}}","last":"{{js .Last}}"{{end}}`
			vars = requestVars{
				First: "Don",
				Last:  "Quixote",
			}
			vars.Type = QueryTypeQuery
			vars.QueryName = "getPerson"
			vars.Args = map[string]string{
				"first": "String!",
				"last":  "String!",
			}
			vars.Returns = []string{
				"first",
				"last",
				"address",
				"age",
			}
			funcs = nil
		})

		It("Parses the template without error", func() {
			_, err := BuildRequestBody(requestTemplate, vars, funcs)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Populates the correct query type", func() {
			buf, _ := BuildRequestBody(requestTemplate, vars, funcs)
			Expect(buf).NotTo(BeNil())
			b, _ := ioutil.ReadAll(buf)
			Expect(b).To(MatchRegexp(`^{"query":"%s`, QueryTypeQuery))
		})

		It("Populates the query argument type spec", func() {
			buf, _ := BuildRequestBody(requestTemplate, vars, funcs)
			Expect(buf).NotTo(BeNil())
			b, _ := ioutil.ReadAll(buf)
			for k, v := range vars.Args {
				Expect(b).To(MatchRegexp(`\$%s: %s`, k, v))
			}
		})

		It("Populates the query method name and arg spec", func() {
			buf, _ := BuildRequestBody(requestTemplate, vars, funcs)
			Expect(buf).NotTo(BeNil())
			b, _ := ioutil.ReadAll(buf)
			for k, _ := range vars.Args {
				Expect(b).To(MatchRegexp(`\\n  %s\([^\)]*%s: \$%s`, vars.QueryName, k, k))
			}
		})

		It("Populates the fields to be returned", func() {
			buf, _ := BuildRequestBody(requestTemplate, vars, funcs)
			Expect(buf).NotTo(BeNil())
			b, _ := ioutil.ReadAll(buf)
			for _, f := range vars.Returns {
				Expect(b).To(MatchRegexp(`\\n    %s`, f))
			}
		})

		It("Processes the request-specific variable template", func() {
			buf, _ := BuildRequestBody(requestTemplate, vars, funcs)
			Expect(buf).NotTo(BeNil())
			b, _ := ioutil.ReadAll(buf)
			v := reflect.ValueOf(vars)
			for k, _ := range vars.Args {
				Expect(b).To(MatchRegexp(`"variables":{[^}]*"%s":"%s"`,
					k, v.FieldByName(strings.Title(k))))
			}
		})

		Context("When additional helper functions are passed in", func() {
			BeforeEach(func() {
				requestTemplate = `{{define "vars"}}"first":"{{js (toUpper .First)}}"{{end}}`
				funcs = template.FuncMap{
					"toUpper": strings.ToUpper,
				}
			})

			It("Merges the supplied function map with the defaults", func() {
				buf, err := BuildRequestBody(requestTemplate, vars, funcs)
				Expect(err).NotTo(HaveOccurred())
				Expect(buf).NotTo(BeNil())
				b, _ := ioutil.ReadAll(buf)
				Expect(b).To(MatchRegexp(strings.ToUpper(vars.First)))
			})
		})

		Context("When the variable template does not escape js for all variables", func() {
			BeforeEach(func() {
				requestTemplate = `{{define "vars"}}"first":"{{js .First}}","last":"{{.Last}}"{{end}}`
			})

			It("Returns nil and an error", func() {
				buf, err := BuildRequestBody(requestTemplate, vars, funcs)
				Expect(buf).To(BeNil())
				Expect(err).To(MatchError(MatchRegexp("All variables must be escaped")))
			})
		})

		Context("When the variable template is not valid", func() {
			BeforeEach(func() {
				requestTemplate = `{{define "vars"}}{{js .First}}`
			})

			It("Returns nil and an error", func() {
				buf, err := BuildRequestBody(requestTemplate, vars, funcs)
				Expect(buf).To(BeNil())
				Expect(err).To(MatchError(MatchRegexp("Unable to parse supplied template")))
			})
		})

		Context("When the template references variables not part of the struct", func() {
			BeforeEach(func() {
				requestTemplate = `{{define "vars"}}"first":"{{js .Foo}}"{{end}}`
			})

			It("Returns nil and an error", func() {
				buf, err := BuildRequestBody(requestTemplate, vars, funcs)
				Expect(buf).To(BeNil())
				Expect(err).To(MatchError(MatchRegexp("Unable to execute template")))
			})
		})
	})
})
