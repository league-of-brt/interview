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

## 1 建议使用第二种，判断err != nil

第一种写法：

```golang
// OneByTeacherIDV1 ...
func (m *Mapper) OneByTeacherIDV1(ctx context.Context, teacherID int64) (*teacher.Schema, error) {
	var schema teacher.Schema
	err := m.DB.WithContext(ctx).Model(&schema).
		Where(m.AliasTag(&m.TeacherID).Eq(), teacherID).
		Select()
	return &schema, err
}
```

`var schema teacher.Schema`定义了一个变量，`&schema`返回其指针，所以这里返回的`*teacher.Schema`永远不为nil！

但是在外部调用方法的时候，我们看到返回一个指针，很有可能就拿这个指针是否为nil去进行逻辑判断了，因为这个指针永不为nil，所以可能会出现问题：

```golang
// DoSomething ...
func DoSomething(ctx context.Context, teacherID int64) error {
	teacherMapper := teacherlogic.New()
	teacher, err := teacherMapper.OneByTeacherIDV1(ctx, teacherID)
	switch {
	case err == nil:
		break
	case pg.IsRecordNotFoundErr(err):
		return nil
	default:
		return err
	}

	// 这里的teacher永远不为nil
	if teacher == nil {
		// doSomething ...
	}

	// 永远走这个逻辑
	if teacher != nil {
		// doSomething ...
	}
	return nil
}
```

正确写法应该是：

```golang
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

这样调用方就没有压力了。

## 2 例外情况

如果返回一个list，这样写行不行呢？

```golang
// ListByTeacherIDsV1 ...
func (m *Mapper) ListByTeacherIDsV1(ctx context.Context, teacherIDs []int64) ([]*teacher.Schema, error) {
	var schemaList []*teacher.Schema
	if len(teacherID) == 0 {
		return schemaList, nil
	}
	err := m.DB.WithContext(ctx).Model(&schemaList).
		Where(m.AliasTag(&m.TeacherID).In(), pg.In(teacherIDs)).
		Select()
	return schemaList, err
}
```

是可以的，因为`var schemaList []*teacher.Schema`没有实际去分配内存空间，要等到`Select()`执行完、有结果之后才会真正赋值。

如果这里找不到结果，`schemaList`为nil。

但是为了代码风格统一，个人建议还是统一判断一下`err != nil`:

```golang
// ListByTeacherIDsV2 ...
func (m *Mapper) ListByTeacherIDsV2(ctx context.Context, teacherIDs []int64) ([]*teacher.Schema, error) {
	var schemaList []*teacher.Schema
	if len(teacherID) == 0 {
		return schemaList, nil
	}
	err := m.DB.WithContext(ctx).Model(&schemaList).
		Where(m.AliasTag(&m.TeacherID).In(), pg.In(teacherIDs)).
		Select()
	if err != nil {
		return nil, err
	}
	return schemaList, err
}
```
