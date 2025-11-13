// import { profile } from "console"


// file này lưu tất cả các đường dẫn api
const API = {
    // Base URLs theo service
    base_vinh: "https://lemarchenoble.id.vn/api/v1",
    base_gateway: "https://lemarchenoble.id.vn/api/v1", // Identity, Profile, Address
    base_agent: "http://localhost:9102/api", // Port 9000 - AI Agent
    base_product: "http://172.26.127.95:9001/v1", // Port 9001 - Product Service
    base_order: "http://172.26.127.95:9002/v1", // Port 9002 - Order Service
    base_transaction: "http://172.26.127.95:9003/v1", // Port 9003 - Transaction Service
    base_analytics: "http://172.26.127.95:9004/v1", // Port 9004 - Analytics & Statistics
    
    // Address Service
    address: {
        provinces: "/address/provinces/get-all",
        districts: "/address/districts/get-all",
        wards: "/address/wards/get-all",
    },
    
    // User & Profile Service
    user: {
        profile: "/profile/users/profiles/get-my-profile",
        addresses: "/profile/users/profiles/profile-subs/get-all-my-sub-profile",
        createAddress: "/profile/users/profiles/profile-subs/insert",
        updateAddress: "/profile/users/profiles/profile-subs/update",
        deleteAddress: "/profile/users/profiles/profile-subs/delete",
        updateProfile: "/profile/users/profiles/update", // {profileId} will be appended
        login: "/identity/auth/login",
        register: "/identity/users/register",
        refresh: "/identity/auth/refresh",
        new_access_token: "/user/new_access_token"
    },
    
    dalogin: {
        ghi: "/dalogin/ghi"
    },
    
    // Category Service
    category:{
        getAll: "/categories/get"
    },
    
    // Product Service
    product:{
        getAll: "/product/getall",
        getdulieu: "data",
        getDetail: "/product/getdetail/"
    },
    
    // Media Service
    media:{
        avtatar: "/media/avatar/",
        product: "/media/products",
    },
    
    // AI Agent Service (Port 9000)
    agent: {
        session: "/session", // POST - Create chat session
        message: "/message", // POST - Send message to agent
    },
    
    // Analytics Service (Port 9004)
    analytics: {
        chatboxReview: "/public/chatbox/review", // POST - Review chatbox response
        customerComplaint: "/public/customer-support/complaint", // POST - Submit complaint
    }

}
export default API