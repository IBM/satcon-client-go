#!/bin/bash

echo "Checking code coverage..."
coverages=($(ginkgo -r -cover client | grep "coverage: " | awk '{print $2}'))

coveragePassedOnce=false
for i in "${coverages[@]}" ; do
  if [ "${i%.*}" -ge "80" ] ; then
    coveragePassedOnce=true
  else
    echo "Code coverage for some package is ${i}. Must be 80%."
    echo "Code coverage not satisfied."
    echo "Run 'ginkgo -r -cover <package>' for more info."
    exit 1
  fi
done

if [ "$coveragePassedOnce" = "false" ] ; then
    echo "Something went wrong running this script. Please file an issue."
    exit 1
fi

echo "Code coverage satisfied!"
