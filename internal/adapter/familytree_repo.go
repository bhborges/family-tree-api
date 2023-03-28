package adapter

import (
	"context"

	"github.com/bhborges/family-tree-api/internal/domain"
)

const qBuildFamilyTreeByPerson = `
	WITH RECURSIVE family_tree AS (
		SELECT id, parent_id, child_id
		FROM relationships
		WHERE child_id = ?
		UNION ALL
		SELECT r.id, r.parent_id, r.child_id
		FROM relationships r
		JOIN family_tree ft ON r.child_id = ft.parent_id
	)
	SELECT DISTINCT p.name as name, COALESCE(p2.name, '') as parent
	FROM family_tree ft
	JOIN people p ON ft.parent_id = p.id
	LEFT JOIN relationships r ON r.child_id = p.id
	LEFT JOIN people p2 ON r.parent_id = p2.id
	UNION ALL
	SELECT DISTINCT p.name as name, COALESCE(p2.name, '') as parent
	FROM people p
	LEFT JOIN relationships r ON r.child_id = p.id
	LEFT JOIN people p2 ON r.parent_id = p2.id
	WHERE p.id = ?`

// BuildFamilyTree builds the family tree of a given person ID, with the person as the root node.
func (pr *PostgresRepository) BuildFamilyTree(ctx context.Context, id string) (*domain.FamilyTree, error) {
	q := pr.db.Raw(qBuildFamilyTreeByPerson, id, id)

	rows, err := q.Rows()
	if err != nil {
		return nil, err
	}

	if rows.Err() != nil {
		return nil, err
	}
	defer rows.Close()

	rs := make(map[string][]string)

	for rows.Next() {
		var name string

		var parent string

		err := rows.Scan(&name, &parent)
		if err != nil {
			return nil, err
		}

		if parent != "" {
			if _, ok := rs[name]; !ok {
				rs[name] = []string{parent}
			} else {
				rs[name] = append(rs[name], parent)
			}
		}
	}

	// create members slice and populate it with Member structs
	ms := make([]*domain.Member, 0)

	for name, parentList := range rs {
		m := &domain.Member{Name: name, Relationships: make([]domain.FamilyRelationship, len(parentList))}
		for i, parent := range parentList {
			m.Relationships[i] = domain.FamilyRelationship{Name: parent, Relationship: "parent"}
		}

		ms = append(ms, m)
	}

	return &domain.FamilyTree{Members: ms}, nil
}
