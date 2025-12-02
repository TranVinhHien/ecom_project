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
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"
import { useTranslations } from "next-intl"
import { useMutation } from "@tanstack/react-query"
import { AxiosError } from "axios"
import API from "@/assets/configs/api"
import { Link, useRouter } from "@/i18n/routing"
import { useToast } from "@/hooks/use-toast"
import { useState } from "react"
import { Loader2 } from "lucide-react"

export default function RegisterForm() {
    const t = useTranslations("Register")
    const router = useRouter()
    const { toast } = useToast()
    const [isLoading, setIsLoading] = useState(false)
    const [errorMessage, setErrorMessage] = useState<string>("")
    const [gender, setGender] = useState<string>("true")
    const [formErrors, setFormErrors] = useState<any>({})

    // Validation functions
    const validateUsername = (username: string) => {
        if (username.length < 8) {
            return "Tên đăng nhập phải có ít nhất 8 ký tự";
        }
        return "";
    };

    const validatePassword = (password: string) => {
        if (password.length < 8) {
            return "Mật khẩu phải có ít nhất 8 ký tự";
        }
        const hasUpperCase = /[A-Z]/.test(password);
        const hasLowerCase = /[a-z]/.test(password);
        const hasNumber = /[0-9]/.test(password);
        const hasSpecialChar = /[!@#$%^&*(),.?":{}|<>]/.test(password);

        if (!hasUpperCase || !hasLowerCase || !hasNumber || !hasSpecialChar) {
            return "Mật khẩu phải có chữ hoa, chữ thường, số và ký tự đặc biệt";
        }
        return "";
    };

    const validateEmail = (email: string) => {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailRegex.test(email)) {
            return "Email không hợp lệ";
        }
        return "";
    };

    const validatePhoneNumber = (phoneNumber: string) => {
        const phoneRegex = /^[0-9]+$/;
        if (!phoneRegex.test(phoneNumber)) {
            return "Số điện thoại chỉ được chứa số";
        }
        if (phoneNumber.length < 10) {
            return "Số điện thoại phải có ít nhất 10 số";
        }
        return "";
    };

    // Register Mutation
    const RegisterMutation = useMutation<any, AxiosError<ResponseType>, any>({
        mutationFn: async (data) => {
            const response = await fetch((API.base_vinh + API.user.register), {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data)
            });
            const result = await response.json();
            
            console.log("=== REGISTER RESPONSE ===");
            console.log("Code:", result.code);
            console.log("Result:", result.result);
            
            if (result.code !== 9997) {
                console.error("REGISTER ERROR:", result);
                throw new Error(result.message || "Đăng ký thất bại");
            }
            
            return result;
        },
    });

    const handleSubmit = async (e: any) => {
        e.preventDefault();
        setIsLoading(true);
        setErrorMessage("");
        setFormErrors({});

        const formData = new FormData(e.target);
        const username = formData.get("username") as string;
        const password = formData.get("password") as string;
        const confirmPassword = formData.get("confirmPassword") as string;
        const email = formData.get("email") as string;
        const phoneNumber = formData.get("phoneNumber") as string;
        const firstName = formData.get("firstName") as string;
        const lastName = formData.get("lastName") as string;
        const dob = formData.get("dob") as string;

        // Validate all fields
        const errors: any = {};
        
        const usernameError = validateUsername(username);
        if (usernameError) errors.username = usernameError;

        const passwordError = validatePassword(password);
        if (passwordError) errors.password = passwordError;

        if (password !== confirmPassword) {
            errors.confirmPassword = "Mật khẩu xác nhận không khớp";
        }

        const emailError = validateEmail(email);
        if (emailError) errors.email = emailError;

        const phoneError = validatePhoneNumber(phoneNumber);
        if (phoneError) errors.phoneNumber = phoneError;

        if (!firstName.trim()) {
            errors.firstName = "Họ không được để trống";
        }

        if (!lastName.trim()) {
            errors.lastName = "Tên không được để trống";
        }

        if (!dob) {
            errors.dob = "Ngày sinh không được để trống";
        }

        // If there are validation errors, show them and stop
        if (Object.keys(errors).length > 0) {
            setFormErrors(errors);
            setErrorMessage("Vui lòng kiểm tra lại thông tin");
            setIsLoading(false);
            return;
        }

        const registerData = {
            username,
            password,
            email,
            firstName,
            lastName,
            dob,
            phoneNumber,
            gender: gender === "true",
            roleCode: "USER"
        };

        console.log("=== REGISTER DATA ===", registerData);

        RegisterMutation.mutate(registerData, {
            onSuccess: (response) => {
                console.log("=== REGISTER SUCCESS ===", response);
                
                toast({
                    title: "Đăng ký thành công!",
                    description: "Bạn có thể đăng nhập ngay bây giờ",
                });

                // Redirect to login page after 1 second
                setTimeout(() => {
                    router.push("/auth/login");
                    setIsLoading(false);
                }, 1000);
            },
            onError: (error: any) => {
                console.error("=== REGISTER ERROR ===", error);
                
                const errorMsg = error.message || "Đăng ký thất bại. Vui lòng thử lại";
                setErrorMessage(errorMsg);
                
                toast({
                    title: "Đăng ký thất bại",
                    description: errorMsg,
                    variant: "destructive",
                });
                setIsLoading(false);
            },
        });
    };

    return (
        <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
            <div className="w-full max-w-2xl">
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
                                <div className="flex flex-col gap-4">
                                    {/* Error Message Display */}
                                    {errorMessage && (
                                        <div className="p-3 rounded-lg bg-red-50 border border-red-200 text-red-800 text-sm">
                                            <p className="font-medium">⚠️ {errorMessage}</p>
                                        </div>
                                    )}
                                    
                                    {/* Username and Password Row */}
                                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                        <div className="grid gap-2">
                                            <Label htmlFor="username">{t("username_label")} <span className="text-red-500">*</span></Label>
                                            <Input
                                                id="username"
                                                name="username"
                                                placeholder="Tối thiểu 8 ký tự"
                                                required
                                                disabled={isLoading}
                                                className={formErrors.username ? "border-red-500" : ""}
                                            />
                                            {formErrors.username && (
                                                <span className="text-xs text-red-500">{formErrors.username}</span>
                                            )}
                                        </div>
                                        <div className="grid gap-2">
                                            <Label htmlFor="email">{t("email_label")} <span className="text-red-500">*</span></Label>
                                            <Input
                                                id="email"
                                                name="email"
                                                type="email"
                                                placeholder="example@email.com"
                                                required
                                                disabled={isLoading}
                                                className={formErrors.email ? "border-red-500" : ""}
                                            />
                                            {formErrors.email && (
                                                <span className="text-xs text-red-500">{formErrors.email}</span>
                                            )}
                                        </div>
                                    </div>

                                    {/* Password Row */}
                                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                        <div className="grid gap-2">
                                            <Label htmlFor="password">{t("password_label")} <span className="text-red-500">*</span></Label>
                                            <Input
                                                id="password"
                                                name="password"
                                                type="password"
                                                placeholder="Tối thiểu 8 ký tự"
                                                required
                                                disabled={isLoading}
                                                className={formErrors.password ? "border-red-500" : ""}
                                            />
                                            {formErrors.password && (
                                                <span className="text-xs text-red-500">{formErrors.password}</span>
                                            )}
                                        </div>
                                        <div className="grid gap-2">
                                            <Label htmlFor="confirmPassword">{t("confirm_password_label")} <span className="text-red-500">*</span></Label>
                                            <Input
                                                id="confirmPassword"
                                                name="confirmPassword"
                                                type="password"
                                                placeholder="Nhập lại mật khẩu"
                                                required
                                                disabled={isLoading}
                                                className={formErrors.confirmPassword ? "border-red-500" : ""}
                                            />
                                            {formErrors.confirmPassword && (
                                                <span className="text-xs text-red-500">{formErrors.confirmPassword}</span>
                                            )}
                                        </div>
                                    </div>

                                    {/* Name Row */}
                                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                        <div className="grid gap-2">
                                            <Label htmlFor="firstName">{t("first_name_label")} <span className="text-red-500">*</span></Label>
                                            <Input
                                                id="firstName"
                                                name="firstName"
                                                placeholder="Họ"
                                                required
                                                disabled={isLoading}
                                                className={formErrors.firstName ? "border-red-500" : ""}
                                            />
                                            {formErrors.firstName && (
                                                <span className="text-xs text-red-500">{formErrors.firstName}</span>
                                            )}
                                        </div>
                                        <div className="grid gap-2">
                                            <Label htmlFor="lastName">{t("last_name_label")} <span className="text-red-500">*</span></Label>
                                            <Input
                                                id="lastName"
                                                name="lastName"
                                                placeholder="Tên"
                                                required
                                                disabled={isLoading}
                                                className={formErrors.lastName ? "border-red-500" : ""}
                                            />
                                            {formErrors.lastName && (
                                                <span className="text-xs text-red-500">{formErrors.lastName}</span>
                                            )}
                                        </div>
                                    </div>

                                    {/* DOB, Phone, Gender Row */}
                                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                                        <div className="grid gap-2">
                                            <Label htmlFor="dob">{t("dob_label")} <span className="text-red-500">*</span></Label>
                                            <Input
                                                id="dob"
                                                name="dob"
                                                type="date"
                                                required
                                                disabled={isLoading}
                                                max={new Date().toISOString().split('T')[0]}
                                                className={formErrors.dob ? "border-red-500" : ""}
                                            />
                                            {formErrors.dob && (
                                                <span className="text-xs text-red-500">{formErrors.dob}</span>
                                            )}
                                        </div>
                                        <div className="grid gap-2">
                                            <Label htmlFor="phoneNumber">{t("phone_label")} <span className="text-red-500">*</span></Label>
                                            <Input
                                                id="phoneNumber"
                                                name="phoneNumber"
                                                type="tel"
                                                placeholder="0123456789"
                                                required
                                                disabled={isLoading}
                                                pattern="[0-9]*"
                                                onKeyPress={(e) => {
                                                    if (!/[0-9]/.test(e.key)) {
                                                        e.preventDefault();
                                                    }
                                                }}
                                                className={formErrors.phoneNumber ? "border-red-500" : ""}
                                            />
                                            {formErrors.phoneNumber && (
                                                <span className="text-xs text-red-500">{formErrors.phoneNumber}</span>
                                            )}
                                        </div>
                                        <div className="grid gap-2">
                                            <Label htmlFor="gender">{t("gender_label")} <span className="text-red-500">*</span></Label>
                                            <Select 
                                                value={gender} 
                                                onValueChange={setGender}
                                                disabled={isLoading}
                                            >
                                                <SelectTrigger>
                                                    <SelectValue placeholder="Chọn giới tính" />
                                                </SelectTrigger>
                                                <SelectContent>
                                                    <SelectItem value="true">{t("male")}</SelectItem>
                                                    <SelectItem value="false">{t("female")}</SelectItem>
                                                </SelectContent>
                                            </Select>
                                        </div>
                                    </div>

                                    <div className="text-xs text-gray-500 mt-2">
                                        <p className="mb-1">* Yêu cầu mật khẩu:</p>
                                        <ul className="list-disc list-inside space-y-1 ml-2">
                                            <li>Tối thiểu 8 ký tự</li>
                                            <li>Có ít nhất 1 chữ hoa (A-Z)</li>
                                            <li>Có ít nhất 1 chữ thường (a-z)</li>
                                            <li>Có ít nhất 1 số (0-9)</li>
                                            <li>Có ít nhất 1 ký tự đặc biệt (!@#$%^&*...)</li>
                                        </ul>
                                    </div>

                                    <Button type="submit" className="w-full mt-2" disabled={isLoading}>
                                        {isLoading ? (
                                            <>
                                                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                                Đang đăng ký...
                                            </>
                                        ) : (
                                            t("register_button")
                                        )}
                                    </Button>
                                </div>
                                <div className="mt-4 text-center text-sm">
                                    {t("login_prompt")}{" "}
                                    <Link href="/auth/login" className="underline underline-offset-4">
                                        {t("login")}
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
