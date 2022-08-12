package myasthurts_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	myasthurts "github.com/jamillosantos/go-my-ast-hurts"
)

var _ = Describe("MethodDescriptor", func() {
	Describe("Compatible", func() {
		It("should find methods compatible", func() {
			ref1 := myasthurts.NewRefType("ref1", nil, nil)
			ref2 := myasthurts.NewRefType("ref2", nil, nil)
			ref3 := myasthurts.NewRefType("ref3", nil, nil)
			md1 := myasthurts.MethodDescriptor{
				BaseType: *myasthurts.NewBaseType(nil, ""),
				Arguments: []myasthurts.MethodArgument{
					{
						Type: ref1,
					},
					{
						Type: ref2,
					},
				},
				Result: []myasthurts.MethodResult{
					{
						Type: ref3,
					},
					{
						Type: ref2,
					},
				},
			}
			md2 := myasthurts.MethodDescriptor{
				BaseType: *myasthurts.NewBaseType(nil, ""),
				Arguments: []myasthurts.MethodArgument{
					{
						Type: ref1,
					},
					{
						Type: ref2,
					},
				},
				Result: []myasthurts.MethodResult{
					{
						Type: ref3,
					},
					{
						Type: ref2,
					},
				},
			}
			Expect(md1.Compatible(&md2)).To(BeTrue())
		})

		It("should find methods not compatible with argument length difference", func() {
			ref1 := myasthurts.NewRefType("ref1", nil, nil)
			ref2 := myasthurts.NewRefType("ref2", nil, nil)
			ref3 := myasthurts.NewRefType("ref3", nil, nil)
			md1 := myasthurts.MethodDescriptor{
				BaseType: *myasthurts.NewBaseType(nil, ""),
				Arguments: []myasthurts.MethodArgument{
					{
						Type: ref1,
					},
					{
						Type: ref2,
					},
				},
				Result: []myasthurts.MethodResult{
					{
						Type: ref3,
					},
					{
						Type: ref2,
					},
				},
			}
			md2 := myasthurts.MethodDescriptor{
				BaseType: *myasthurts.NewBaseType(nil, ""),
				Arguments: []myasthurts.MethodArgument{
					{
						Type: ref1,
					},
				},
				Result: []myasthurts.MethodResult{
					{
						Type: ref3,
					},
					{
						Type: ref2,
					},
				},
			}
			Expect(md1.Compatible(&md2)).To(BeFalse())
		})

		It("should find methods not compatible with result length difference", func() {
			ref1 := myasthurts.NewRefType("ref1", nil, nil)
			ref2 := myasthurts.NewRefType("ref2", nil, nil)
			ref3 := myasthurts.NewRefType("ref3", nil, nil)
			md1 := myasthurts.MethodDescriptor{
				BaseType: *myasthurts.NewBaseType(nil, ""),
				Arguments: []myasthurts.MethodArgument{
					{
						Type: ref1,
					},
					{
						Type: ref2,
					},
				},
				Result: []myasthurts.MethodResult{
					{
						Type: ref3,
					},
					{
						Type: ref2,
					},
				},
			}
			md2 := myasthurts.MethodDescriptor{
				BaseType: *myasthurts.NewBaseType(nil, ""),
				Arguments: []myasthurts.MethodArgument{
					{
						Type: ref1,
					},
					{
						Type: ref2,
					},
				},
				Result: []myasthurts.MethodResult{
					{
						Type: ref3,
					},
				},
			}
			Expect(md1.Compatible(&md2)).To(BeFalse())
		})

		It("should find methods not compatible with argument types different", func() {
			ref1 := myasthurts.NewRefType("ref1", nil, nil)
			ref2 := myasthurts.NewRefType("ref2", nil, nil)
			ref3 := myasthurts.NewRefType("ref3", nil, nil)
			md1 := myasthurts.MethodDescriptor{
				BaseType: *myasthurts.NewBaseType(nil, ""),
				Arguments: []myasthurts.MethodArgument{
					{
						Type: ref1,
					},
					{
						Type: ref3,
					},
				},
				Result: []myasthurts.MethodResult{
					{
						Type: ref3,
					},
					{
						Type: ref2,
					},
				},
			}
			md2 := myasthurts.MethodDescriptor{
				BaseType: *myasthurts.NewBaseType(nil, ""),
				Arguments: []myasthurts.MethodArgument{
					{
						Type: ref1,
					},
					{
						Type: ref2,
					},
				},
				Result: []myasthurts.MethodResult{
					{
						Type: ref3,
					},
					{
						Type: ref2,
					},
				},
			}
			Expect(md1.Compatible(&md2)).To(BeFalse())
		})

		It("should find methods not compatible with result types different", func() {
			ref1 := myasthurts.NewRefType("ref1", nil, nil)
			ref2 := myasthurts.NewRefType("ref2", nil, nil)
			ref3 := myasthurts.NewRefType("ref3", nil, nil)
			md1 := myasthurts.MethodDescriptor{
				BaseType: *myasthurts.NewBaseType(nil, ""),
				Arguments: []myasthurts.MethodArgument{
					{
						Type: ref1,
					},
					{
						Type: ref2,
					},
				},
				Result: []myasthurts.MethodResult{
					{
						Type: ref3,
					},
					{
						Type: ref1,
					},
				},
			}
			md2 := myasthurts.MethodDescriptor{
				BaseType: *myasthurts.NewBaseType(nil, ""),
				Arguments: []myasthurts.MethodArgument{
					{
						Type: ref1,
					},
					{
						Type: ref2,
					},
				},
				Result: []myasthurts.MethodResult{
					{
						Type: ref3,
					},
					{
						Type: ref2,
					},
				},
			}
			Expect(md1.Compatible(&md2)).To(BeFalse())
		})
	})
})
