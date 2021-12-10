from django.urls import path
from . import views

urlpatterns = [
    path('', views.index, name="main"),
    path('all-products', views.products, name="products"),
    path('product/<str:product_id>', views.product, name="product"),
    path('reviews', views.reviews, name="reviews"),
    path('contact', views.contact, name="contact"),
    path('logIn', views.logIn, name="logIn"), 
    path('dashboard', views.dashBoard, name="admin-dashboard"), 
    path('add-product',views.addproduct, name='add-product'),
    path('edit-product/<str:product_id>',views.editproduct, name='edit-product'),
    path('edit-manage/<str:product_id>',views.editStatus, name='edit-manage'), # use to edit
    path('add-manage/<str:product_id>',views.addStatus, name='add-manage'), # use to add
    path('del-manage/<str:product_id>',views.delStatus, name='del-manage'), # use to delete
    path('add-comment',views.addComment, name='add-comment'), # use to add comment
    path('li-manage',views.logInManage, name='login-manage'),
    path('checkout-manage/<str:product_id>',views.checkOut, name='checkout-manage')
]
