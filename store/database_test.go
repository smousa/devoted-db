package store_test

import (
	. "github.com/smousa/devoted-db/store"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Database", func() {

	var db *Database

	BeforeEach(func() {
		db = NewDatabase()
	})

	Context("with a populated database", func() {
		BeforeEach(func() {
			db.Set("foo", "bar")
		})

		It("should overwrite the value with the same value", func() {
			db.Set("foo", "bar")
			value, ok := db.Get("foo")
			Ω(ok).Should(BeTrue())
			Ω(value).Should(Equal("bar"))
			Ω(db.Count("bar")).Should(Equal(1))
		})

		It("should overwrite the value with a different value", func() {
			db.Set("foo", "baz")
			value, ok := db.Get("foo")
			Ω(ok).Should(BeTrue())
			Ω(value).Should(Equal("baz"))
			Ω(db.Count("bar")).Should(BeZero())
			Ω(db.Count("baz")).Should(Equal(1))
		})

		It("should get the value", func() {
			value, ok := db.Get("foo")
			Ω(ok).Should(BeTrue())
			Ω(value).Should(Equal("bar"))
		})

		It("should count the value", func() {
			Ω(db.Count("bar")).Should(Equal(1))
		})

		It("should delete the key", func() {
			db.Delete("foo")
			value, ok := db.Get("foo")
			Ω(ok).Should(BeFalse())
			Ω(value).Should(BeEmpty())
			Ω(db.Count("bar")).Should(BeZero())
		})

		It("should commit", func() {
			Ω(db.Commit()).Should(Equal(db))
		})

		It("should not rollback", func() {
			s, err := db.Rollback()
			Ω(err).Should(Equal(ErrTxNotFound))
			Ω(s).Should(BeNil())
		})

		Context("with a transaction", func() {
			var tx Store
			BeforeEach(func() {
				tx = db.Begin()
			})

			It("should set a key and commit", func() {
				tx.Set("foo", "baz")
				value, ok := tx.Get("foo")
				Ω(ok).Should(BeTrue())
				Ω(value).Should(Equal("baz"))
				Ω(tx.Count("bar")).Should(BeZero())
				Ω(tx.Count("baz")).Should(Equal(1))

				value, ok = db.Get("foo")
				Ω(ok).Should(BeTrue())
				Ω(value).Should(Equal("bar"))
				Ω(db.Count("bar")).Should(Equal(1))
				Ω(db.Count("baz")).Should(BeZero())

				s := tx.Commit()
				Ω(s).Should(Equal(db))

				value, ok = db.Get("foo")
				Ω(ok).Should(BeTrue())
				Ω(value).Should(Equal("baz"))
				Ω(db.Count("bar")).Should(BeZero())
				Ω(db.Count("baz")).Should(Equal(1))
			})

			It("should set a key and rollback", func() {
				tx.Set("foo", "baz")
				value, ok := tx.Get("foo")
				Ω(ok).Should(BeTrue())
				Ω(value).Should(Equal("baz"))
				Ω(tx.Count("bar")).Should(BeZero())
				Ω(tx.Count("baz")).Should(Equal(1))

				value, ok = db.Get("foo")
				Ω(ok).Should(BeTrue())
				Ω(value).Should(Equal("bar"))
				Ω(db.Count("bar")).Should(Equal(1))
				Ω(db.Count("baz")).Should(BeZero())

				s, err := tx.Rollback()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(s).Should(Equal(db))

				value, ok = db.Get("foo")
				Ω(ok).Should(BeTrue())
				Ω(value).Should(Equal("bar"))
				Ω(db.Count("bar")).Should(Equal(1))
				Ω(db.Count("baz")).Should(BeZero())
			})

			It("should delete a key and commit", func() {
				tx.Delete("foo")
				value, ok := tx.Get("foo")
				Ω(ok).Should(BeFalse())
				Ω(value).Should(BeEmpty())
				Ω(tx.Count("bar")).Should(BeZero())

				value, ok = db.Get("foo")
				Ω(ok).Should(BeTrue())
				Ω(value).Should(Equal("bar"))
				Ω(db.Count("bar")).Should(Equal(1))

				s := tx.Commit()
				Ω(s).Should(Equal(db))

				value, ok = db.Get("foo")
				Ω(ok).Should(BeFalse())
				Ω(value).Should(BeEmpty())
				Ω(db.Count("bar")).Should(BeZero())
			})

			It("should delete a key and rollback", func() {
				tx.Delete("foo")
				value, ok := tx.Get("foo")
				Ω(ok).Should(BeFalse())
				Ω(value).Should(BeEmpty())
				Ω(tx.Count("bar")).Should(BeZero())

				value, ok = db.Get("foo")
				Ω(ok).Should(BeTrue())
				Ω(value).Should(Equal("bar"))
				Ω(db.Count("bar")).Should(Equal(1))

				s, err := tx.Rollback()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(s).Should(Equal(db))

				value, ok = db.Get("foo")
				Ω(ok).Should(BeTrue())
				Ω(value).Should(Equal("bar"))
				Ω(db.Count("bar")).Should(Equal(1))
			})

		})
	})

	Context("with an empty database", func() {
		It("should set a key", func() {
			db.Set("foo", "bar")
			value, ok := db.Get("foo")
			Ω(ok).Should(BeTrue())
			Ω(value).Should(Equal("bar"))
			Ω(db.Count("bar")).Should(Equal(1))
		})

		It("should not find a key", func() {
			value, ok := db.Get("foo")
			Ω(ok).Should(BeFalse())
			Ω(value).Should(BeEmpty())
		})

		It("should not count a value", func() {
			Ω(db.Count("bar")).Should(BeZero())
		})

		It("should commit", func() {
			Ω(db.Commit()).Should(Equal(db))
		})

		It("should not rollback", func() {
			s, err := db.Rollback()
			Ω(err).Should(Equal(ErrTxNotFound))
			Ω(s).Should(BeNil())
		})

		Context("with a transaction", func() {
			var tx Store
			BeforeEach(func() {
				tx = db.Begin()
			})

			It("should set a key and commit", func() {
				tx.Set("foo", "bar")
				value, ok := tx.Get("foo")
				Ω(ok).Should(BeTrue())
				Ω(value).Should(Equal("bar"))
				Ω(tx.Count("bar")).Should(Equal(1))

				value, ok = db.Get("foo")
				Ω(ok).Should(BeFalse())
				Ω(value).Should(BeEmpty())
				Ω(db.Count("bar")).Should(BeZero())

				s := tx.Commit()
				Ω(s).Should(Equal(db))

				value, ok = db.Get("foo")
				Ω(ok).Should(BeTrue())
				Ω(value).Should(Equal("bar"))
				Ω(db.Count("bar")).Should(Equal(1))
			})

			It("should set a key and rollback", func() {
				tx.Set("foo", "bar")
				value, ok := tx.Get("foo")
				Ω(ok).Should(BeTrue())
				Ω(value).Should(Equal("bar"))
				Ω(tx.Count("bar")).Should(Equal(1))

				value, ok = db.Get("foo")
				Ω(ok).Should(BeFalse())
				Ω(value).Should(BeEmpty())
				Ω(db.Count("bar")).Should(BeZero())

				s, err := tx.Rollback()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(s).Should(Equal(db))

				value, ok = db.Get("foo")
				Ω(ok).Should(BeFalse())
				Ω(value).Should(BeEmpty())
				Ω(db.Count("bar")).Should(BeZero())
			})
		})
	})
})
