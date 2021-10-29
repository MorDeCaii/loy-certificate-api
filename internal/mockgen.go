package internal

//go:generate mockgen -destination=./mocks/repo_mock.go -package=mocks github.com/Mordecaii/loy-certificate-api/internal/app/repo EventRepo
//go:generate mockgen -destination=./mocks/sender_mock.go -package=mocks github.com/Mordecaii/loy-certificate-api/internal/app/sender EventSender
