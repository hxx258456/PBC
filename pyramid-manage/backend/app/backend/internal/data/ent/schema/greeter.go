package schema

import "entgo.io/ent"

// Greeter holds the schema definition for the Greeter entity.
type Greeter struct {
	ent.Schema
}

// Fields of the Greeter.
func (Greeter) Fields() []ent.Field {
	return nil
}

// Edges of the Greeter.
func (Greeter) Edges() []ent.Edge {
	return nil
}
