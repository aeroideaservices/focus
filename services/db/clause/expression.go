package clause

import "gorm.io/gorm/clause"

// Btw
type Btw struct {
	Column interface{}
	First  interface{}
	Second interface{}
}

func (btw Btw) Build(builder clause.Builder) {
	builder.WriteQuoted(btw.Column)
	builder.WriteString(" BETWEEN ")
	builder.AddVar(builder, btw.First)
	builder.WriteString(" AND ")
	builder.AddVar(builder, btw.Second)
}

func (btw Btw) NegationBuild(builder clause.Builder) {
	builder.WriteQuoted(btw.Column)
	builder.WriteString(" NOT BETWEEN ")
	builder.AddVar(builder, btw.First)
	builder.WriteString(" AND ")
	builder.AddVar(builder, btw.Second)
}

// Ilike whether string matches regular expression
type Ilike clause.Like

func (like Ilike) Build(builder clause.Builder) {
	builder.AddVar(builder, like.Column)
	builder.WriteString(" ILIKE ")
	builder.AddVar(builder, like.Value)
}

func (like Ilike) NegationBuild(builder clause.Builder) {
	builder.AddVar(builder, like.Column)
	builder.WriteString(" NOT ILIKE ")
	builder.AddVar(builder, like.Value)
}

// Cast преобразование к типу
type Cast struct {
	Column any
	Type   string
}

func (cast Cast) Build(builder clause.Builder) {
	builder.WriteString("CAST(")
	builder.WriteQuoted(cast.Column)
	builder.WriteString(" AS ")
	builder.WriteString(cast.Type)
	builder.WriteString(")")
}

type Case struct {
	WhenThen []WhenThen
	Else     any
}

type WhenThen struct {
	When any
	Then any
}

func (cas Case) Build(builder clause.Builder) {
	builder.WriteString("CASE")
	for _, wt := range cas.WhenThen {
		builder.WriteString(" WHEN ")
		builder.AddVar(builder, wt.When)
		builder.WriteString(" THEN ")
		builder.AddVar(builder, wt.Then)
	}
	if cas.Else != nil {
		builder.WriteString(" ELSE ")
		builder.AddVar(builder, cas.Else)
	}

	builder.WriteString(" END")
}
