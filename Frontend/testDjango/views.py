from django.http.response import HttpResponseNotAllowed, HttpResponseRedirect
from django.shortcuts import render
import requests
import json
import random
import os
from dotenv import load_dotenv

load_dotenv()
# Create your views here.
def index(request):
    return render(request, 'index.html')

def products(request):
    res = requests.get(
        os.getenv("URL")+"/api/products/list",
        headers={'Authorization': 'Bearer '+os.getenv("TOKEN")}
    )
    if res.status_code != 204:
        context = {'products': res.json()["result"]}
    else:
        context = {'products': None}
    return render(request, 'products.html', context)

def product(request, product_id):
    res = requests.get(
        os.getenv("URL")+"/api/products/"+product_id,
        headers={'Authorization': 'Bearer '+os.getenv("TOKEN")}
    )
    context = {'product': res.json()["result"]}
    return render(request, 'product.html', context)

def reviews(request):
    res = requests.get(
        os.getenv("URL")+"/api/reviews/list",
        headers={'Authorization': 'Bearer '+os.getenv("TOKEN")}
    )
    context = {'reviews': res.json()["result"]}
    return render(request, 'reviews.html', context)

def addComment(request):
    if request.method == "POST":
        payload = {
            'review_id': format(random.randint(0000,9999), '04d'),
            'reviewer': request.POST['reviewer'],
            'rating': int(request.POST['rating']),
            'comment': request.POST['comment']
        }
        res = requests.post(
            os.getenv("URL")+"/api/reviews/add",
            headers={'Authorization': 'Bearer '+os.getenv("TOKEN")},
            data=json.dumps(payload)
        )
        if res.status_code == 200:
            return HttpResponseRedirect('/reviews')
        else:
            return render(request, 'failedPage.html', {'failedMessage': 'Add Comment failed', 'reason': res.json()})

def contact(request):
    return render(request, 'contact.html')

def logIn(request):
    return render(request, 'logIn.html')

def logInManage(request):
    if request.method == 'POST':
        payload = {
            "username": request.POST['username'],
            "password": request.POST['password']
        }
        res = requests.post(
            os.getenv("URL")+"/api/admin/login",
            headers={'Authorization': 'Bearer '+os.getenv("TOKEN")},
            data=json.dumps(payload)
        )
        if res.status_code == 200:
            return HttpResponseRedirect("/dashboard")
        else:
            return render(request, 'failedPage.html', {'failedMessage': 'Log In failed', 'reason': res.json()})

def dashBoard(request):
    res = requests.get(
        os.getenv("URL")+"/api/products/list",
        headers={'Authorization': 'Bearer '+os.getenv("TOKEN")}
    )
    if res.status_code != 204:
        context = {'products': res.json()["result"]}
    else:
        context = {'products': None}
    return render(request, 'admin_dashboard.html', context)

def addproduct(request):
    prod_id = format(random.randint(0000,9999), '04d')
    return render(request, 'add_product.html', {"product_id": prod_id})

def addStatus(request, product_id):
    if request.method != 'POST':
        return render(request, 'failedPage.html', {'failedMessage': 'Access Denied', 'reason': 'authorized user only!'})
    payload = {
            "prod_id": product_id,
            "prod_name": request.POST['prod_name'],
            "prod_detail": request.POST['prod_detail'],
            "prod_price": int(request.POST['prod_price']),
            "prod_quantity": int(request.POST['prod_quantity'])
    }
    addRes = requests.post(
        os.getenv("URL")+"/api/products/add",
        headers={'Content-Type': "application/json", 'Accept': "application/json", 'Authorization': 'Bearer '+os.getenv("TOKEN")},
        data=json.dumps(payload)
    )
    print(addRes.json())
    for key, file in request.FILES.items():
        res = requests.post(
            os.getenv("URL")+"/api/products/upload-image",
            headers={'Authorization': 'Bearer '+os.getenv("TOKEN"), 'prod_id': product_id},
            files={key: file.read()}
        )
        print(res.json())
        if res.status_code == 200:
            return HttpResponseRedirect("/dashboard")
        else:
            return render(request, 'failedPage.html', {'failedMessage': 'Add failed', 'reason': res.json()})
    return HttpResponseRedirect("/dashboard")

def editproduct(request, product_id):
    res = requests.get(
        os.getenv("URL")+"/api/products/"+product_id,
        headers={'Authorization': 'Bearer '+os.getenv("TOKEN")}
    )
    context = {'product': res.json()["result"]}
    return render(request, 'edit_product.html', context)

def editStatus(request, product_id):
    if request.method != 'POST':
        return render(request, 'failedPage.html', {'failedMessage': 'Access Denied', 'reason': 'authorized user only!'})

    payload = {
            "prod_id": product_id,
            "prod_name": request.POST['prod_name'],
            "prod_detail": request.POST['prod_detail'],
            "prod_price": int(request.POST['prod_price']),
            "prod_quantity": int(request.POST['prod_quantity'])
    }
    editRes = requests.put(
        os.getenv("URL")+"/api/products/edit",
        headers={'Content-Type': "application/json", 'Accept': "application/json", 'Authorization': 'Bearer '+os.getenv("TOKEN")},
        data=json.dumps(payload)
    )
    print(editRes.json())
    if request.FILES is None:
        return HttpResponseRedirect("/dashboard")
    for key, file in request.FILES.items():
        res = requests.post(
            os.getenv("URL")+"/api/products/upload-image",
            headers={'Authorization': 'Bearer '+os.getenv("TOKEN"), 'prod_id': product_id},
            files={key: file.read()}
        )
        print(res.json())
        if res.status_code == 200:
            return HttpResponseRedirect("/dashboard")
        else:
            return render(request, 'failedPage.html', {'failedMessage': 'Edit failed', 'reason': res.json()})
    return HttpResponseRedirect("/dashboard")

def delStatus(request, product_id):
    if request.method == "POST":
        delRes = requests.delete(
            os.getenv("URL")+"/api/products/del",
            headers={'Authorization': 'Bearer '+os.getenv("TOKEN")},
            data=json.dumps({'prod_id': product_id})
        )
        print(delRes.json())
        res = requests.post(
            os.getenv("URL")+"/api/products/delete-image",
            headers={'Authorization': 'Bearer '+os.getenv("TOKEN"), 'prod_id': product_id},
        )
        print(res.json())
        if res.status_code == 200:
            return HttpResponseRedirect("/dashboard")
        else:
            return render(request, 'failedPage.html', {'failedMessage': 'Delete failed', 'reason': res.json()})
        
    return render(request, 'failedPage.html', {'failedMessage': 'Delete failed'})


def checkOut(request, product_id):
    if request.method == "POST":
        payload = {
            'prod_id': product_id,
            'prod_name': request.POST['name'],
            'prod_detail': request.POST['detail'],
            'prod_price': int(request.POST['price']),
            'prod_quantity': int(request.POST['total_quan']) - int(request.POST['number'])
        }
        # print(payload)
        editRes = requests.put(
            os.getenv("URL")+"/api/products/edit",
            headers={'Content-Type': "application/json", 'Accept': "application/json", 'Authorization': 'Bearer '+os.getenv("TOKEN")},
            data=json.dumps(payload)
        )
        print(editRes.json())
        if editRes.status_code == 200:
            return render(request, 'successPage.html', {'successMessage': 'ชำระเงินเสร็จสิ้น ขอบคุณที่ใช้บริการ ขอเวลาอีกไม่นาน เราจะส่งของให้คุณ'})
        else:
            return render(request, 'failedPage.html', {'failedMessage': 'Checkout failed', 'reason': editRes.json()})
    return render(request, 'failedPage.html', {'failedMessage': 'Checkout failed'})