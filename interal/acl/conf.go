package acl

const conf = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && regexMatch(r.obj, p.obj) && regexMatch(r.act, p.act)
`
