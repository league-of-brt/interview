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

## 1 建议使用第二种，判断err != nil

第一种写法：

```golang
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
```

`var schema organizationuserteachersnapshot.Schema`定义了一个变量，`&schema`返回其指针，所以这里返回的`*organizationuserteachersnapshot.Schema`永远不为nil！

但是在外部调用方法的时候，我们看到返回一个指针，很有可能就拿这个指针是否为nil去进行逻辑判断了，因为这个指针永不为nil，所以可能会出现问题：

![pic](https://brt-1303999354.cos.ap-shanghai.myqcloud.com/%E5%BE%AE%E4%BF%A1%E5%9B%BE%E7%89%87_20210308112125.png)

正确写法应该是：

```golang
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

这样调用方就没有调用压力了。

## 2 例外情况

如果返回一个list，这样写行不行呢？

```golang
// ListByOrgIDNTeacherIDV3 ...
func (m *Mapper) ListByOrgIDNTeacherIDV3(ctx context.Context, orgID, teacherID, version int64) ([]*organizationuserteachersnapshot.Schema, error) {
	var schemaList []*organizationuserteachersnapshot.Schema
	err := m.DB.WithContext(ctx).Model(&schemaList).
		Where(m.AliasTag(&m.OrganizationID).Eq(), orgID).
		Where(m.AliasTag(&m.TeacherID).Eq(), teacherID).
		Where(m.AliasTag(&m.VersionID).Eq(), version).
		Select()
	return schemaList, err
}
```

是可以的，因为`var schemaList []*organizationuserteachersnapshot.Schema`并没有实际去分配内存空间，要等到`Select()`执行完、有结果之后才会真正赋值。如果这里找不到结果，`schemaList`为nil。

但是为了代码风格统一，个人建议还是统一判断一下`err != nil`:

```golang
// ListByOrgIDNTeacherIDV4 ...
func (m *Mapper) ListByOrgIDNTeacherIDV4(ctx context.Context, orgID, teacherID, version int64) ([]*organizationuserteachersnapshot.Schema, error) {
	var schemaList []*organizationuserteachersnapshot.Schema
	err := m.DB.WithContext(ctx).Model(&schemaList).
		Where(m.AliasTag(&m.OrganizationID).Eq(), orgID).
		Where(m.AliasTag(&m.TeacherID).Eq(), teacherID).
		Where(m.AliasTag(&m.VersionID).Eq(), version).
		Select()
	if err != nil {
		return nil, err
	}
	return schemaList, err
}
```
