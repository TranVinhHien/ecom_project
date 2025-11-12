"use client"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { useTranslations } from "next-intl"
import { useMutation } from "@tanstack/react-query"
import { AxiosError } from "axios"
import API from "@/assets/configs/api"
import { jwtDecode } from "jwt-decode"
import { ACCESS_TOKEN, INFO_USER } from "@/assets/configs/request"
import { cookies } from "@/assets/helpers"
import { Link, useRouter } from "@/i18n/routing"
import { useToast } from "@/hooks/use-toast"
import { useState } from "react"
import { Loader2 } from "lucide-react"
import { tokenRefreshService } from "@/lib/tokenRefreshService"

export default function LoginForm() {
    const t = useTranslations("Login")
    const router = useRouter()
    const { toast } = useToast()
    const [isLoading, setIsLoading] = useState(false)
    const [errorMessage, setErrorMessage] = useState<string>("")

    // Mutation 1: Login
    const LoginMutation = useMutation<any, AxiosError<ResponseType>, any>({
        mutationFn: async (data) => {
            const response = await fetch((API.base_vinh + API.user.login), {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    username: data.username,
                    password: data.password
                })
            });
            const result = await response.json();
            
            console.log("=== LOGIN RESPONSE ===");
            console.log("Code:", result.code);
            console.log("Token:", result.result?.token);
            if (result.result?.token == null){
                // throw new Error("Không nhận được token từ server");
            }
            if (result.code !== 10000) {
                console.error("LOGIN ERROR:", result);
                throw new Error("Tài khoảng không chính xác");
            }
            if (!result.result?.token) {
                throw new Error("Mật khẩu không chính xác");
            }
            return result;
        },
    });

    // Mutation 2: Get User Profile
    const getProfileMutation = useMutation<any, AxiosError<ResponseType>, any>({
        mutationFn: async (token: string) => {
            const response = await fetch(API.base_vinh + API.user.profile, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`,
                },
            });
            const result = await response.json();
            
            console.log("=== USER PROFILE RESPONSE ===");
            console.log("Code:", result.code);
            console.log("Profile:", result.result);
            
            if (result.code !== 10000) {
                console.error("PROFILE ERROR:", result);
                throw new Error("Không thể lấy thông tin người dùng");
            }
            
            return result.result;
        },
    });

    // Mutation 3: Get User Addresses
    const getAddressesMutation = useMutation<any, AxiosError<ResponseType>, any>({
        mutationFn: async (token: string) => {
            const response = await fetch(API.base_vinh + API.user.addresses, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`,
                },
            });
            const result = await response.json();
            
            console.log("=== USER ADDRESSES RESPONSE ===");
            console.log("Code:", result.code);
            console.log("Addresses:", result.result);
            
            if (result.code !== 10000) {
                console.error("ADDRESSES ERROR:", result);
                throw new Error("Không thể lấy thông tin địa chỉ");
            }
            
            return result.result;
        },
    });

    const handleSubmit = async (e: any) => {
        e.preventDefault();
        setIsLoading(true);
        setErrorMessage(""); // Clear previous errors

        const formData = new FormData(e.target);
        const loginData = {
            username: formData.get("username"),
            password: formData.get("password"),
        };

        // Step 1: Login
        LoginMutation.mutate(loginData, {
            onSuccess: async (loginResponse) => {
                const token = loginResponse.result?.token;
                
                if (!token) {
                    toast({
                        title: "Lỗi đăng nhập",
                        description: "Không nhận được token từ server",
                        variant: "destructive",
                    });
                    setIsLoading(false);
                    return;
                }

                try {
                    // Decode token to get expiry
                    const decoded: any = jwtDecode(token);
                    
                    // Step 2: Get User Profile
                    const profileData = await getProfileMutation.mutateAsync(token);
                    
                    // Step 3: Get User Addresses
                    const addressesData = await getAddressesMutation.mutateAsync(token);
                    
                    // Format addresses for localStorage
                    const formattedAddresses = addressesData.map((addr: any) => ({
                        id_address: addr.id,
                        address: `${addr.address.other}, ${addr.address.ward.fullName}, ${addr.address.ward.district.fullName}, ${addr.address.ward.district.province.fullName}`,
                        phone_number: addr.phoneNumber,
                        name: addr.name,
                        ward: addr.address.ward,
                        district: addr.address.ward.district,
                        province: addr.address.ward.district.province,
                    }));

                    // Prepare user data
                    const userData = {
                        id: profileData.id,
                        userId: profileData.userId,
                        name: `${profileData.firstName} ${profileData.lastName}`,
                        firstName: profileData.firstName,
                        lastName: profileData.lastName,
                        dob: profileData.dob,
                        phone_number: profileData.phoneNumber,
                        gender: profileData.gender ? "Nam" : "Nữ",
                        addresses: formattedAddresses,
                    };

                    // Save to cookies and localStorage
                    // cookies.setCookieValues(ACCESS_TOKEN, token, decoded?.exp);
                    cookies.setCookieValues(ACCESS_TOKEN, token, decoded?.exp+300);
                    localStorage.setItem(INFO_USER, JSON.stringify(userData));

                    console.log("=== SAVED USER DATA ===");
                    console.log("Token saved to cookies");
                    console.log("User data:", userData);

                    // Show success toast
                    toast({
                        title: "Đăng nhập thành công!",
                        description: `Chào mừng ${userData.name}`,
                    });

                    // Initialize token refresh service after successful login
                    console.log("🔄 Initializing token refresh service...");
                    tokenRefreshService.initialize();

                    // Redirect to home
                    setTimeout(() => {
                        router.push("/");
                        setIsLoading(false);
                    }, 500);

                } catch (error: any) {
                    console.error("=== ERROR DURING LOGIN FLOW ===", error);
                    
                    const errorMsg = error.message || "Có lỗi xảy ra trong quá trình đăng nhập";
                    setErrorMessage(errorMsg);
                    
                    toast({
                        title: "Lỗi",
                        description: errorMsg,
                        variant: "destructive",
                    });
                    setIsLoading(false);
                }
            },
            onError: (error: any) => {
                console.error("=== LOGIN ERROR ===", error);
                
                const errorMsg = error.message || "Tên đăng nhập hoặc mật khẩu không đúng";
                setErrorMessage(errorMsg);
                
                toast({
                    title: "Đăng nhập thất bại",
                    description: errorMsg,
                    variant: "destructive",
                });
                setIsLoading(false);
            },
        });
    };

    return (
        <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
            <div className="w-full max-w-sm">
                <div className={cn("flex flex-col gap-6")}>
                    <Card>
                        <CardHeader>
                            <CardTitle className="text-2xl">{t("title")}</CardTitle>
                            <CardDescription>
                                {t("description")}
                            </CardDescription>
                        </CardHeader>
                        <CardContent>
                            <form onSubmit={handleSubmit}>
                                <div className="flex flex-col gap-6">
                                    {/* Error Message Display */}
                                    {errorMessage && (
                                        <div className="p-3 rounded-lg bg-red-50 border border-red-200 text-red-800 text-sm">
                                            <p className="font-medium">⚠️ {errorMessage}</p>
                                        </div>
                                    )}
                                    
                                    <div className="grid gap-2">
                                        <Label htmlFor="username"> {t("username_label")}</Label>
                                        <Input
                                            id="username"
                                            name="username"
                                            placeholder="hienlazada"
                                            required
                                            disabled={isLoading}
                                        />
                                    </div>
                                    <div className="grid gap-2">
                                        <div className="flex items-center">
                                            <Label htmlFor="password"> {t("password_label")}</Label>
                                            <a
                                                href="#"
                                                className="ml-auto inline-block text-sm underline-offset-4 hover:underline"
                                            >
                                                {t("forgot_password")}
                                            </a>
                                        </div>
                                        <Input 
                                            autoComplete="current-password" 
                                            id="password" 
                                            name="password" 
                                            type="password" 
                                            required 
                                            disabled={isLoading}
                                        />
                                    </div>
                                    <Button type="submit" className="w-full" disabled={isLoading}>
                                        {isLoading ? (
                                            <>
                                                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                                Đang đăng nhập...
                                            </>
                                        ) : (
                                            t("login_button")
                                        )}
                                    </Button>
                                    <Button variant="outline" className="w-full" disabled={isLoading}>
                                        {t("login_with_google")}
                                    </Button>
                                </div>
                                <div className="mt-4 text-center text-sm">
                                    {t("sign_up_prompt")}{" "}
                                    <Link href="/auth/register" className="underline underline-offset-4">
                                        {t("sign_up")}
                                    </Link>
                                </div>
                            </form>
                        </CardContent>
                    </Card>
                </div>
            </div>
        </div>
    )
}
