package example1_test

import (
	"testing"

	"github.com/r-erema/go_sendbox/learning/software_design_patterns/behavioral/observer/example1"
	"github.com/stretchr/testify/assert"
)

func TestObserver(t *testing.T) {
	t.Parallel()

	item := example1.NewItem("iPhone 13")
	customer := example1.NewCustomer("1", "M. Salah")
	customer2 := example1.NewCustomer("2", "T.A. Arnold")
	customer3 := example1.NewCustomer("3", "J. Henderson")

	item.Register(customer)
	item.Register(customer2)
	item.Register(customer3)
	item.Deregister(customer2)
	item.SetAsAvailable()

	assert.Equal(t, "I got your message about `iPhone 13 is available`, thanks! Best regards, M. Salah", customer.Log()[0])
	assert.Empty(t, customer2.Log())
	assert.Equal(t, "I got your message about `iPhone 13 is available`, thanks! Best regards, J. Henderson", customer3.Log()[0])
}
