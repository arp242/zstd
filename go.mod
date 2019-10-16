module zgo.at/utils

go 1.12

require (
	github.com/pkg/errors v0.8.0
	github.com/teamwork/test v0.0.0-20190410143529-8897d82f8d46
	github.com/teamwork/utils v0.0.0-00010101000000-000000000000 // indirect
)

replace github.com/teamwork/utils => ./ // Because test depends on Teamwork utils.
