import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSidenavModule } from '@angular/material/sidenav';
import { RouterOutlet } from '@angular/router';
import { CartComponent } from './components/cart/cart.component';
import { ProductListComponent } from './components/product-list/product-list.component';
import { CartService } from './services/cart/cart.service';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    RouterOutlet,
    ProductListComponent,
    CartComponent,
    MatButtonModule,
    MatSidenavModule,
    MatInputModule,
    FormsModule,
    MatIconModule,
    MatCardModule,
  ],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss',
})
export class AppComponent {
  title = 'shop';
  user = '';

  constructor(private _cartService: CartService) {}

  getCart(): void {
    console.log('Getting cart');
    this._cartService.getCart(this.user).subscribe();
  }
}
