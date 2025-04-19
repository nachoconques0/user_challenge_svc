
mockgen --source=internal/coupon/coupon.go --destination=internal/mocks/mock_coupon.go --package=mocks --mock_names=Repository=MockCouponRepository,Service=MockCouponService,Server=MockCouponServer
mockgen --source=internal/shopping_cart/shopping_cart.go --destination=internal/mocks/mock_shopping_cart.go --package=mocks --mock_names=Repository=MockShoppingCartRepository,Service=MockShoppingCartService,Server=MockShoppingCartServer
