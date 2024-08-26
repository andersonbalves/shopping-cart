import { CommonModule, CurrencyPipe } from '@angular/common';
import { Component, Input, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatListModule } from '@angular/material/list';
import { Observable } from 'rxjs';
import { CartItemModel } from '../../model/cart-item.model';
import { CartService } from '../../services/cart/cart.service';

@Component({
  selector: 'app-cart',
  standalone: true,
  imports: [
    CurrencyPipe,
    CommonModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    FormsModule,
    MatInputModule,
    MatFormFieldModule,
    MatListModule,
  ],
  templateUrl: './cart.component.html',
  styleUrls: ['./cart.component.scss'],
})
export class CartComponent implements OnInit {
  cart$: Observable<CartItemModel[]>;
  total$: Observable<number>;
  _user = '';

  @Input() set user(user: string) {
    this._user = user;
  }

  constructor(private _cartService: CartService) {
    this.cart$ = this._cartService.cart$;
    this.total$ = this._cartService.getTotal();
  }

  ngOnInit(): void {
    this._cartService.getCart(this._user).subscribe();
  }

  plusQuantity(product: CartItemModel): void {
    const updatedProduct = { ...product, quantity: product.quantity + 1 };
    this._cartService.updateCartItem(this._user, updatedProduct).subscribe();
  }

  minusQuantity(product: CartItemModel): void {
    if (product.quantity > 1) {
      const updatedProduct = { ...product, quantity: product.quantity - 1 };
      this._cartService.updateCartItem(this._user, updatedProduct).subscribe();
    } else {
      this.removeItem(product.id);
    }
  }

  removeItem(productId: number): void {
    this._cartService.removeFromCart(this._user, productId).subscribe();
  }
}
