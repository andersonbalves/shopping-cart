import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { BehaviorSubject, map, Observable, tap } from 'rxjs';
import { CartItemModel } from '../../model/cart-item.model';
import { ProductModel } from '../../model/product.model';

@Injectable({
  providedIn: 'root',
})
export class CartService {
  private baseUrl = '/restapis/oglczdbcsj/test/_user_request_';
  private cartSubject = new BehaviorSubject<CartItemModel[]>([]);
  cart$ = this.cartSubject.asObservable();

  constructor(private http: HttpClient) {}

  getCart(user: string): Observable<CartItemModel[]> {
    return this.http
      .get<any[]>(`${this.baseUrl}/shopping-cart`, {
        params: { userId: user },
      })
      .pipe(
        tap((response) => {
          this.updateCart(response);
        })
      );
  }

  addToCart(user: string, product: ProductModel): Observable<CartItemModel[]> {
    return this.http
      .post<any[]>(`${this.baseUrl}/shopping-cart`, {
        userId: user,
        productId: String(product.id),
        productName: product.name,
        quantity: 1,
        unitPrice: product.price,
      })
      .pipe(
        tap((response) => {
          this.updateCart(response);
        })
      );
  }

  updateCartItem(
    user: string,
    product: CartItemModel
  ): Observable<CartItemModel[]> {
    return this.http
      .put<any[]>(`${this.baseUrl}/shopping-cart`, {
        userId: user,
        productId: String(product.id),
        quantity: product.quantity,
      })
      .pipe(
        tap((response) => {
          this.updateCart(response);
        })
      );
  }

  removeFromCart(user: string, productId: number): Observable<CartItemModel[]> {
    return this.http
      .delete<any[]>(`${this.baseUrl}/shopping-cart`, {
        params: { userId: user, productId: String(productId) },
      })
      .pipe(
        tap((response) => {
          this.updateCart(response);
        })
      );
  }

  private updateCart(response: any[]): void {
    const cart = response.map((item) => ({
      id: item.productId,
      name: item.productName,
      price: item.unitPrice,
      quantity: item.quantity,
    }));
    this.cartSubject.next([...cart]); // Emite um novo array de produtos
  }

  getTotal(): Observable<number> {
    return this.cart$.pipe(
      map((cart) =>
        cart.reduce((sum, item) => sum + item.price * item.quantity, 0)
      )
    );
  }
}
