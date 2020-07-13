package actions_test

import (
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
			"org_id":    "String!",
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
				"$org_id: String!",
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
				"org_id: $org_id",
				"flavor: $flavor",
				"dimension: $dimension",
			))

			Expect(argVarList).NotTo(HaveSuffix(", "))
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
			requestTemplate = `{{define "vars"}}"first":"{{.First}}","last":"{{.Last}}"{{end}}`
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
	})
})
