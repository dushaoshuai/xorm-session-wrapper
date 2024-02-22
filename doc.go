/*
Package wrapper is a thin wrapper of *xorm.Session.
It aims at eliminating tedious if statements.

List can be replaced by List2:

	type ListReq struct {
		CondA  []int64
		CondB  []string
		CondC  string
		CondD  bool
		CondE  int64
		CondF  *Ranger
		Limit  int
		Offset int
	}

	type Model struct{}

	func List(ctx context.Context, req *ListReq) (int64, []*Model, error) {
		var sess *xorm.Session
		sess.Context(ctx)

		if len(req.CondA) != 0 {
			sess.Where("column_a IN ?", req.CondA)
		}
		if len(req.CondB) != 0 {
			sess.Where("column_b IN ?", req.CondB)
		}
		if condC := strings.TrimSpace(req.CondC); condC != "" {
			sess.Where("column_c LIKE ?", "%"+condC+"%")
		}

		sess.Where("column_d = ?", req.CondD)

		if req.CondE != 0 {
			sess.Where("column_e = ?", req.CondE)
		}
		if req.CondF != nil {
			sess.Where("column_f BETWEEN ? AND ?",
				req.CondF.Start, req.CondF.End)
		}

		var models []*Model
		count, err := sess.
			Desc("id").
			Limit(req.Limit, req.Offset).
			FindAndCount(&models)

		if err != nil {
			return 0, nil, err
		}
		return count, models, nil
	}

	func List2(ctx context.Context, req *ListReq) (int64, []*Model, error) {
		sess := NewSession((*xorm.Session)(nil))

		var models []*Model
		count, err := sess.Context(ctx).
			In("column_a", req.CondA).
			In("column_b", req.CondB).
			Like("column_c", req.CondC).
			Where("column_d = ?", req.CondD).
			Equal("column_e", req.CondE).
			Between("column_f", req.CondF).
			Desc("id").
			Limit(req.Limit, req.Offset).
			FindAndCount(&models)

		if err != nil {
			return 0, nil, err
		}
		return count, models, nil
	}
*/
package wrapper
