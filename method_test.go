package myasthurts

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MethodDescriptor", func() {
	Describe("Compatible", func() {
		It("should find methods compatible", func() {
			ref1 := NewRefType("ref1", nil, nil)
			ref2 := NewRefType("ref2", nil, nil)
			ref3 := NewRefType("ref3", nil, nil)
			md1 := MethodDescriptor{
				baseType: *NewBaseType(nil, ""),
				Arguments: []MethodArgument{
					{
						Type: ref1,
					},
					{
						Type: ref2,
					},
				},
				Result: []MethodResult{
					{
						Type: ref3,
					},
					{
						Type: ref2,
					},
				},
			}
			md2 := MethodDescriptor{
				baseType: *NewBaseType(nil, ""),
				Arguments: []MethodArgument{
					{
						Type: ref1,
					},
					{
						Type: ref2,
					},
				},
				Result: []MethodResult{
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
			ref1 := NewRefType("ref1", nil, nil)
			ref2 := NewRefType("ref2", nil, nil)
			ref3 := NewRefType("ref3", nil, nil)
			md1 := MethodDescriptor{
				baseType: *NewBaseType(nil, ""),
				Arguments: []MethodArgument{
					{
						Type: ref1,
					},
					{
						Type: ref2,
					},
				},
				Result: []MethodResult{
					{
						Type: ref3,
					},
					{
						Type: ref2,
					},
				},
			}
			md2 := MethodDescriptor{
				baseType: *NewBaseType(nil, ""),
				Arguments: []MethodArgument{
					{
						Type: ref1,
					},
				},
				Result: []MethodResult{
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
			ref1 := NewRefType("ref1", nil, nil)
			ref2 := NewRefType("ref2", nil, nil)
			ref3 := NewRefType("ref3", nil, nil)
			md1 := MethodDescriptor{
				baseType: *NewBaseType(nil, ""),
				Arguments: []MethodArgument{
					{
						Type: ref1,
					},
					{
						Type: ref2,
					},
				},
				Result: []MethodResult{
					{
						Type: ref3,
					},
					{
						Type: ref2,
					},
				},
			}
			md2 := MethodDescriptor{
				baseType: *NewBaseType(nil, ""),
				Arguments: []MethodArgument{
					{
						Type: ref1,
					},
					{
						Type: ref2,
					},
				},
				Result: []MethodResult{
					{
						Type: ref3,
					},
				},
			}
			Expect(md1.Compatible(&md2)).To(BeFalse())
		})

		It("should find methods not compatible with argument types different", func() {
			ref1 := NewRefType("ref1", nil, nil)
			ref2 := NewRefType("ref2", nil, nil)
			ref3 := NewRefType("ref3", nil, nil)
			md1 := MethodDescriptor{
				baseType: *NewBaseType(nil, ""),
				Arguments: []MethodArgument{
					{
						Type: ref1,
					},
					{
						Type: ref3,
					},
				},
				Result: []MethodResult{
					{
						Type: ref3,
					},
					{
						Type: ref2,
					},
				},
			}
			md2 := MethodDescriptor{
				baseType: *NewBaseType(nil, ""),
				Arguments: []MethodArgument{
					{
						Type: ref1,
					},
					{
						Type: ref2,
					},
				},
				Result: []MethodResult{
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
			ref1 := NewRefType("ref1", nil, nil)
			ref2 := NewRefType("ref2", nil, nil)
			ref3 := NewRefType("ref3", nil, nil)
			md1 := MethodDescriptor{
				baseType: *NewBaseType(nil, ""),
				Arguments: []MethodArgument{
					{
						Type: ref1,
					},
					{
						Type: ref2,
					},
				},
				Result: []MethodResult{
					{
						Type: ref3,
					},
					{
						Type: ref1,
					},
				},
			}
			md2 := MethodDescriptor{
				baseType: *NewBaseType(nil, ""),
				Arguments: []MethodArgument{
					{
						Type: ref1,
					},
					{
						Type: ref2,
					},
				},
				Result: []MethodResult{
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
