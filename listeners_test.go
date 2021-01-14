package myasthurts_test

import (
	"errors"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	myasthurts "github.com/jamillosantos/go-my-ast-hurts"
)

type (
	beforeListener func(*myasthurts.ParsePackageContext, string) error
	afterListener  func(*myasthurts.ParsePackageContext, string, error) error
)

func (listener beforeListener) BeforeFile(ctx *myasthurts.ParsePackageContext, filePath string) error {
	return listener(ctx, filePath)
}

func (listener afterListener) AfterFile(ctx *myasthurts.ParsePackageContext, filePath string, err error) error {
	return listener(ctx, filePath, err)
}

var _ = Describe("Listeners", func() {
	Describe("BeforeFile", func() {
		It("should call before each file", func() {
			files := make([]string, 0)
			env, err := myasthurts.NewEnvironmentWithListener(beforeListener(func(ctx *myasthurts.ParsePackageContext, filePath string) error {
				files = append(files, filePath)
				return nil
			}))
			Expect(err).ToNot(HaveOccurred())
			pkg, err := env.ParseDir("./data/parse_dir")
			Expect(err).ToNot(HaveOccurred())
			Expect(files).To(HaveLen(2))
			Expect(files).To(ConsistOf("data/parse_dir/home.go", "data/parse_dir/user.go"))
			Expect(pkg.Structs).To(HaveLen(2))
		})

		It("should skip after the first file", func() {
			env, err := myasthurts.NewEnvironmentWithListener(beforeListener(func(ctx *myasthurts.ParsePackageContext, filePath string) error {
				if strings.HasSuffix(filePath, "home.go") {
					return myasthurts.Skip
				}
				return nil
			}))
			Expect(err).ToNot(HaveOccurred())
			pkg, err := env.ParseDir("./data/parse_dir")
			Expect(err).ToNot(HaveOccurred())
			Expect(pkg.Structs).To(HaveLen(1))
			Expect(pkg.Structs[0].Name()).To(Equal("User"))
		})
	})

	Describe("AfterFile", func() {
		It("should call after each file", func() {
			files := make([]string, 0)
			var (
				env *myasthurts.Environment
				err error
			)
			env, err = myasthurts.NewEnvironmentWithListener(afterListener(func(ctx *myasthurts.ParsePackageContext, filePath string, err error) error {
				Expect(err).ToNot(HaveOccurred())
				files = append(files, filePath)
				return nil
			}))
			Expect(err).ToNot(HaveOccurred())
			pkg, err := env.ParseDir("./data/parse_dir")
			Expect(err).ToNot(HaveOccurred())
			Expect(files).To(HaveLen(2))
			Expect(files).To(ConsistOf("data/parse_dir/home.go", "data/parse_dir/user.go"))
			Expect(pkg.Structs).To(HaveLen(2))
		})

		It("should abort when error", func() {
			env, err := myasthurts.NewEnvironmentWithListener(beforeListener(func(ctx *myasthurts.ParsePackageContext, filePath string) error {
				if strings.HasSuffix(filePath, "home.go") {
					return errors.New("forced error")
				}
				return nil
			}))
			Expect(err).ToNot(HaveOccurred())
			pkg, err := env.ParseDir("./data/parse_dir")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("forced error"))
			Expect(pkg).To(BeNil())
		})
	})
})
