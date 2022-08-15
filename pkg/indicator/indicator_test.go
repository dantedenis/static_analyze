package indicator

import (
	"github.com/stretchr/testify/assert"
	"log"
	"static-analyze/pkg/indicator/mock"
	"testing"
)

var (
	pairStr = []string{"VALUE1", "VALUE3", "", "TEST"}
)

func TestNew(t *testing.T) {
	c := New(mock.Proto{}, pairStr)
	assert.NotNil(t, c)

	assert.Nil(t, c.Updater())

	_, _, err := c.Get("VALUE1", 5)
	assert.NotNil(t, err)

	get, t2, err := c.Get("VALUE1", 60)
	assert.Nil(t, err)
	assert.NotNil(t, get, t2)
	log.Println(get, t2)
}
