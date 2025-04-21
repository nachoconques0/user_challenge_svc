echo "🔄 Generating mocks..."

mockgen --source=pkg/challenge/internal/aggregate/user/user.go --destination=pkg/challenge/internal/mocks/mock_user_aggregate.go --package=mocks --mock_names=Aggregate=MockUserAggregate
mockgen --source=pkg/challenge/internal/service/user/service.go --destination=pkg/challenge/internal/mocks/mock_user_service.go --package=mocks --mock_names=Service=MockUserService
mockgen --source=pkg/challenge/pubsub/publisher.go --destination=pkg/challenge/internal/mocks/mock_publisher.go --package=mocks --mock_names=Publisher=MockPublisher

echo "✅ Mocks generated!"