#!/bin/bash

set -e  # Exit on error

# Set up coverage directory
COVERAGE_DIR="coverage"
mkdir -p $COVERAGE_DIR

echo "Running tests with coverage..."
if ! go test -coverprofile=$COVERAGE_DIR/coverage.out -covermode=atomic ./...; then
    echo "Error: Tests failed"
    exit 1
fi

echo "Generating HTML coverage report..."
if ! go tool cover -html=$COVERAGE_DIR/coverage.out -o $COVERAGE_DIR/coverage.html; then
    echo "Error: Failed to generate HTML coverage report"
    exit 1
fi

echo "Generating coverage summary..."
if ! go tool cover -func=$COVERAGE_DIR/coverage.out > $COVERAGE_DIR/coverage.txt; then
    echo "Error: Failed to generate coverage summary"
    exit 1
fi

echo "Coverage Summary:"
cat $COVERAGE_DIR/coverage.txt

# Check if coverage meets minimum threshold
MIN_COVERAGE=80
COVERAGE=$(go tool cover -func=$COVERAGE_DIR/coverage.out | grep total | awk '{print $3}' | sed 's/%//')

if (( $(echo "$COVERAGE < $MIN_COVERAGE" | bc -l) )); then
    echo "Error: Test coverage ($COVERAGE%) is below minimum threshold ($MIN_COVERAGE%)"
    exit 1
else
    echo "Test coverage ($COVERAGE%) meets minimum threshold ($MIN_COVERAGE%)"
fi

# Generate coverage badge
echo "Generating coverage badge..."
BADGE_COLOR="brightgreen"
if (( $(echo "$COVERAGE < 60" | bc -l) )); then
    BADGE_COLOR="red"
elif (( $(echo "$COVERAGE < 80" | bc -l) )); then
    BADGE_COLOR="yellow"
fi

cat > $COVERAGE_DIR/badge.svg << EOF
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
    <path fill="#$BADGE_COLOR" d="M67 0h32v20H67z"/>
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

echo "Coverage report generated in $COVERAGE_DIR/"
echo "HTML report: $COVERAGE_DIR/coverage.html"
echo "Text summary: $COVERAGE_DIR/coverage.txt"
echo "Badge: $COVERAGE_DIR/badge.svg" 