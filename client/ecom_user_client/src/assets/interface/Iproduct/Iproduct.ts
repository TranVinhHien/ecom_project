interface DescriptionAttr {
    attr_id: string;
    name: string;
    value: string;
}

interface Rating {
    comment: string;
    create_date: string;
    name: string;
    star: number;
}

interface SkuAttr {
    image: Valid<string>;
    name: string;
    product_sku_attr_id: string;
    value: string;
}

interface Sku {
    price: number;
    product_sku_id: string;
    sku_stock: number;
    value: string;
}

interface Spu {
    average_star: string;
    brand_id: string;
    category_id: string;
    description: string;
    image: string;
    key: string;
    media: string;
    name: string;
    price: number;
    products_spu_id: string;
    short_description: string;
    total_rating: number;
}

interface ProductData {
    description_attrs: DescriptionAttr[];
    ratings: Rating[];
    sku_attrs: SkuAttr[];
    sku: Sku[];
    spu: Spu;
}
interface ProductSimple {
    products_spu_id: string;
    name: string;
    image: string;
    price: number;
    discount: number;
    average_star: number;
    total_rating: number; 
    brand_id:string;
    key: string;
}
