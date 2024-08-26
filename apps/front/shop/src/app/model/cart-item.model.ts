import { ProductModel } from './product.model';

export interface CartItemModel extends ProductModel {
  quantity: number;
}
