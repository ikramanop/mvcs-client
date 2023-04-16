package utility

import (
	"bufio"
	"strings"

	"github.com/ikramanop/mvcs-client/app/model"
)

func ReadStdin(reader *bufio.Reader) (string, error) {
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", model.WrapError(model.StdinError, err)
	}
	input = strings.TrimSpace(input)

	return input, nil
}
