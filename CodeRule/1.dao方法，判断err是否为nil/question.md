# dao方法，判断err是否为nil，讨论两种写法的区别？

```golang
// OneByTeacherIDV1 ...
func (m *Mapper) OneByTeacherIDV1(ctx context.Context, teacherID int64) (*teacher.Schema, error) {
	var schema teacher.Schema
	err := m.DB.WithContext(ctx).Model(&schema).
		Where(m.AliasTag(&m.TeacherID).Eq(), teacherID).
		Select()
	return &schema, err
}

// OneByTeacherIDV2 ...
func (m *Mapper) OneByTeacherIDV2(ctx context.Context, teacherID int64) (*teacher.Schema, error) {
	var schema teacher.Schema
	err := m.DB.WithContext(ctx).Model(&schema).
		Where(m.AliasTag(&m.TeacherID).Eq(), teacherID).
		Select()
	if err != nil {
		return nil, err
	}
	return &schema, err
}
```