package shortener

// Service will own the business rules of the product:
//   - validate the long URL
//   - generate a short code
//   - persist via store.Store
//   - look up a code for redirects
//
// Why this package exists: handlers should parse HTTP and return status
// codes. They should not generate IDs or decide collision policy.
//
// Code generation (hash vs counter vs random + Base62) is a system-design
// choice we will pick deliberately in the next step, not hide in a handler.
type Service struct {
	// Store and code generator will be injected here later.
	// Empty for now so the package compiles and the dependency direction is clear.
}
