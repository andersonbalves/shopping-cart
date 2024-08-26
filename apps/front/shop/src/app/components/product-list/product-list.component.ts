import { CommonModule, CurrencyPipe } from '@angular/common';
import { Component, Input, OnInit } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { ProductModel } from '../../model/product.model';
import { CartService } from '../../services/cart/cart.service';
import { ProductsService } from '../../services/products/products.service';

@Component({
  selector: 'app-product-list',
  imports: [
    CurrencyPipe,
    CommonModule,
    MatCardModule,
    MatButtonModule,
    MatIconModule,
  ],
  templateUrl: './product-list.component.html',
  styleUrls: ['./product-list.component.scss'],
  standalone: true,
})
export class ProductListComponent implements OnInit {
  products: ProductModel[] = [];
  _user = '';

  @Input() set user(user: string) {
    this._user = user;
  }

  constructor(
    private _productsService: ProductsService,
    private _cartService: CartService
  ) {}

  ngOnInit(): void {
    this._productsService.getProducts().subscribe((data) => {
      this.products = data;
    });
  }

  addToCart(product: ProductModel): void {
    this._cartService.addToCart(this._user, product).subscribe();
  }
}
