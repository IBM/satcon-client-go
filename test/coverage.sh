#!/bin/bash

echo "Checking code coverage..."
coverages=($(ginkgo -r -cover client/actions | grep "coverage: " | awk '{print $2}'))

coveragePassed=true
for i in "${coverages[@]}" ; do
  if [ "${i%.*}" -lt "80" ] ; then
    echo "Code coverage for some package is ${i}. Must be 80%."
    coveragePassed=false
  fi
done

if [ "$coveragePassed" = false ] ; then
    echo "Code coverage not satisfied."
    echo "Run 'ginkgo -r -cover <package>' for more info."
    exit 1
fi

echo "Code coverage satisfied!"
