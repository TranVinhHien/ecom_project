interface Role {
    id: string;
    name: string;
    code: string;
}
interface Address{
    address:string;
    id_address:string;
    phone_number:string;
}

interface UserLoginType {
    gender: string;
    dob: string;
    email: string;
    name:string;
    image:Valid<string>;
    addess:Address[];
}