package topic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTopic(t *testing.T) {
	tp := NewTopic("User", "datainq")
	assert.NotNil(t, tp)
}

func TestTopicBroadcastToOne(t *testing.T) {
	tp := NewTopic("User", "datainq")
	go tp.Run()
	c := tp.GetPipe()
	l, canc := tp.AddListener()
	ex := Post{Author: "datainq", Message: "Test content"}
	c <- ex
	m := <-l
	canc()
	assert.Equal(t, ex, m)
}

func TestTopicBroadcastToMany(t *testing.T) {
	tp := NewTopic("User", "datainq")
	go tp.Run()
	c := tp.GetPipe()
	l0, canc0 := tp.AddListener()
	l1, canc1 := tp.AddListener()
	ex := Post{Author: "datainq", Message: "Test content"}
	c <- ex
	close(c)
	m := <-l0
	assert.Equal(t, ex, m)
	m = <-l1
	assert.Equal(t, ex, m)
	canc0()
	canc1()
}

func TestTopicForwardMany(t *testing.T) {
	tp := NewTopic("User", "datainq")
	go tp.Run()
	c0 := tp.GetPipe()
	c1 := tp.GetPipe()

	l0, canc0 := tp.AddListener()
	ex := Post{Author: "datainq", Message: "Test content"}
	c0 <- ex
	m := <-l0
	assert.Equal(t, ex, m)

	ex = Post{Author: "datainq", Message: "Test content1"}
	c1 <- ex
	close(c1)
	m = <-l0
	assert.Equal(t, ex, m)
	canc0()
}
