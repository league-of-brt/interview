# 讨论一下dao方法这两种写法的区别？

```golang
package daoreturn

// OneByOrgIDNTeacherIDV1 ...
func (m *Mapper) OneByOrgIDNTeacherIDV1(ctx context.Context, orgID, teacherID, version int64) (*organizationuserteachersnapshot.Schema, error) {
	var schema organizationuserteachersnapshot.Schema
	err := m.DB.WithContext(ctx).Model(&schema).
		Where(m.AliasTag(&m.OrganizationID).Eq(), orgID).
		Where(m.AliasTag(&m.TeacherID).Eq(), teacherID).
		Where(m.AliasTag(&m.VersionID).Eq(), version).
		Select()
	return &schema, err
}

// OneByOrgIDNTeacherIDV2 ...
func (m *Mapper) OneByOrgIDNTeacherIDV2(ctx context.Context, orgID, teacherID, version int64) (*organizationuserteachersnapshot.Schema, error) {
	var schema organizationuserteachersnapshot.Schema
	err := m.DB.WithContext(ctx).Model(&schema).
		Where(m.AliasTag(&m.OrganizationID).Eq(), orgID).
		Where(m.AliasTag(&m.TeacherID).Eq(), teacherID).
		Where(m.AliasTag(&m.VersionID).Eq(), version).
		Select()
	if err != nil {
		return nil, err
	}
	return &schema, err
}
```