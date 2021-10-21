package internal_test

import (
	. "github.com/homedepot/go-clouddriver/internal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Map", func() {

	var m, actual map[string]interface{}

	Describe("#DeleteNilValues", func() {

		JustBeforeEach(func() {
			actual = DeleteNilValues(m)
		})

		When("the map is empty", func() {
			BeforeEach(func() {
				m = map[string]interface{}{}
			})

			It("returns an empty map", func() {
				Expect(actual).To(Equal(m))
			})
		})

		When("the map has no nil-valued keys", func() {
			BeforeEach(func() {
				m = map[string]interface{}{
					"a": 1,
					"b": true,
					"c": "three",
				}
			})

			It("returns the map unchanged", func() {
				Expect(actual).To(Equal(m))
			})
		})

		When("the map has nil-valued keys", func() {
			BeforeEach(func() {
				m = map[string]interface{}{
					"a": 1,
					"b": nil,
					"c": "three",
				}
			})

			It("succeeds", func() {
				expected := map[string]interface{}{
					"a": 1,
					"c": "three",
				}
				Expect(actual).To(Equal(expected))
			})
		})

		When("the map has nested nil-valued keys", func() {
			BeforeEach(func() {
				m = map[string]interface{}{
					"a": 1,
					"b": "nil",
					"c": "three",
					"d": map[string]interface{}{
						"e": nil,
						"f": "f",
						"g": map[string]interface{}{
							"h": nil,
						},
					},
				}
			})

			It("succeeds", func() {
				expected := map[string]interface{}{
					"a": 1,
					"b": "nil",
					"c": "three",
					"d": map[string]interface{}{
						"f": "f",
						"g": map[string]interface{}{},
					},
				}
				Expect(actual).To(Equal(expected))
			})
		})
	})
})
