# Put this script in the in the folder below the go /src folder before running it.

export GOPATH=`PWD`
export GOBIN=`PWD`/bin

echo "GET"

# Testing library.
go get -d "gopkg.in/check.v1"
[ $? -ne 0 ] && exit 1

# Genetic library.
go get -d "github.com/glemzurg/go-genetic"
[ $? -ne 0 ] && exit 1

echo "FMT" ; go fmt github.com/glemzurg...
[ $? -ne 0 ] && exit 1

echo "TEST" ; go test github.com/glemzurg... -p=1
[ $? -ne 0 ] && exit 1

echo "INSTALL" ; go install github.com/glemzurg...
[ $? -ne 0 ] && exit 1

# Run an experiment.
echo "RUN" ; bin/crazy_card_hyper -scorer=src/github.com/glemzurg/go-genetic/examples/experiments/cmd/crazy_card_hyper/config/scorer.json \
                                  -genetic=src/github.com/glemzurg/go-genetic/examples/experiments/cmd/crazy_card_hyper/config/genetic.json \
                                  -sorter=src/github.com/glemzurg/go-genetic/examples/experiments/cmd/crazy_card_hyper/config/sorter_hypervolume.json
                                  -selector=src/github.com/glemzurg/go-genetic/examples/experiments/cmd/crazy_card_hyper/config/truncate_selector.json
