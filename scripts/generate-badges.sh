#!/bin/bash

set -e  # Exit on error

echo "Generating badges..."

# Get coverage percentage
if [ ! -f "coverage/coverage.out" ]; then
    echo "Error: Coverage file not found. Run coverage.sh first."
    exit 1
fi

COVERAGE=$(go tool cover -func=coverage/coverage.out | grep total | awk '{print $3}' | sed 's/%//')
if [ -z "$COVERAGE" ]; then
    echo "Error: Could not determine coverage percentage"
    exit 1
fi

echo "Coverage: $COVERAGE%"

# Generate coverage badge
echo "Generating coverage badge..."
cat > coverage/badge.svg << EOF
<svg xmlns="http://www.w3.org/2000/svg" width="99" height="20">
  <linearGradient id="b" x2="0" y2="100%">
    <stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
    <stop offset="1" stop-opacity=".1"/>
  </linearGradient>
  <mask id="a">
    <rect width="99" height="20" rx="3" fill="#fff"/>
  </mask>
  <g mask="url(#a)">
    <path fill="#555" d="M0 0h67v20H0z"/>
    <path fill="#4c1" d="M67 0h32v20H67z"/>
    <path fill="url(#b)" d="M0 0h99v20H0z"/>
  </g>
  <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11">
    <text x="33.5" y="15" fill="#010101" fill-opacity=".3">coverage</text>
    <text x="33.5" y="14">coverage</text>
    <text x="83" y="15" fill="#010101" fill-opacity=".3">$COVERAGE%</text>
    <text x="83" y="14">$COVERAGE%</text>
  </g>
</svg>
EOF

# Update README with new badge
if [ ! -f "README.md" ]; then
    echo "Error: README.md not found"
    exit 1
fi

echo "Updating README badges..."
sed -i "s/coverage-[0-9]*%/coverage-$COVERAGE%/" README.md

# Generate benchmark badge
echo "Running benchmarks..."
BENCHMARKS=$(go test -bench=. -benchmem ./... 2>&1 | grep -c "PASS")
if [ $? -ne 0 ]; then
    echo "Error: Benchmarks failed"
    exit 1
fi

if [ $BENCHMARKS -gt 0 ]; then
    echo "Benchmarks passed"
    sed -i "s/benchmarks-[a-z]*/benchmarks-passing/" README.md
else
    echo "Benchmarks failed"
    sed -i "s/benchmarks-[a-z]*/benchmarks-failing/" README.md
fi

# Generate load test badge
echo "Running load tests..."
LOAD_TESTS=$(go test -v ./tests/load/... 2>&1 | grep -c "PASS")
if [ $? -ne 0 ]; then
    echo "Error: Load tests failed"
    exit 1
fi

if [ $LOAD_TESTS -gt 0 ]; then
    echo "Load tests passed"
    sed -i "s/load%20tests-[a-z]*/load%20tests-passing/" README.md
else
    echo "Load tests failed"
    sed -i "s/load%20tests-[a-z]*/load%20tests-failing/" README.md
fi

echo "Badge generation completed" 