
#!/bin/bash

set -e

echo "Running unit tests..."
# TODO catch dataraces
# ginkgo -r -p -keepGoing -trace -randomizeAllSpecs -progress --race cli client
ginkgo -r -p -keepGoing -trace -progress cli client
