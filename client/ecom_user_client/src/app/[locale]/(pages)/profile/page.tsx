"use client";

import { useEffect, useState } from "react";
import { useRouter } from "@/i18n/routing";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { User, Phone, Calendar, MapPin, Edit, Loader2, Plus, Trash2, Save, X } from "lucide-react";
import { getCookieValues } from "@/assets/helpers/cookies";
import { ACCESS_TOKEN, INFO_USER } from "@/assets/configs/request";
import ROUTER from "@/assets/configs/routers";
import { useToast } from "@/hooks/use-toast";
import AddressDialog from "@/components/AddressDialog";
import type { UserAddress } from "@/types/address.types";
import API from "@/assets/configs/api";
import { apiClient } from "@/lib/apiClient";
import { UserProfile } from "@/types/user.types";


export default function ProfilePage() {
  const router = useRouter();
  const { toast } = useToast();
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [addresses, setAddresses] = useState<UserAddress[]>([]);
  const [isLoadingAddresses, setIsLoadingAddresses] = useState(false);
  
  const [addressDialogOpen, setAddressDialogOpen] = useState(false);
  const [editingAddress, setEditingAddress] = useState<UserAddress | undefined>();
  
  // Edit profile states
  const [isEditing, setIsEditing] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const [editedProfile, setEditedProfile] = useState<Partial<UserProfile>>({});




  useEffect(() => {
    const token = getCookieValues<string>(ACCESS_TOKEN);
    
    if (!token) {
      toast({
        title: "Vui lòng đăng nhập",
        description: "Bạn cần đăng nhập để xem thông tin cá nhân",
        variant: "destructive",
      });
      
      setTimeout(() => {
        router.push(ROUTER.auth.login);
      }, 1500);
      return;
    }

    const userInfo = localStorage.getItem(INFO_USER);
    if (userInfo) {
      try {
        const userData = JSON.parse(userInfo);
        setProfile(userData);
      } catch (error) {
        console.error("Error parsing user data:", error);
      }
    }
    
    loadAddresses();
    fetchLatestUserProfile();
  }, [router, toast]);

  const fetchLatestUserProfile = async () => {
    try {
      const token = getCookieValues<string>(ACCESS_TOKEN);
      if (!token) return;

      // Get profile
      const profileResponse = await apiClient.get(API.user.profile);
      const profileData = profileResponse.data.result;

      // Get addresses
      const addressesResponse = await apiClient.get(API.user.addresses);
      const addressesData = addressesResponse.data.result;

      // Format addresses
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

      // Save to localStorage
      localStorage.setItem(INFO_USER, JSON.stringify(userData));
      setProfile(userData);
      
      toast({
        title: "Thành công",
        description: "Đã cập nhật thông tin mới nhất",
      });

      console.log("✅ Profile updated from API:", userData);
    } catch (error) {
      console.error("Error fetching profile:", error);
    } finally {
      setIsLoading(false);
    }
  };

  const loadAddresses = async () => {
    setIsLoadingAddresses(true);
    try {
      const token = getCookieValues<string>(ACCESS_TOKEN);
      const response = await fetch(API.base_vinh + API.user.addresses, {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      });

      const data = await response.json();
      if (data.code === 10000) {
        setAddresses(data.result || []);
      }
    } catch (error) {
      console.error("Error loading addresses:", error);
    } finally {
      setIsLoadingAddresses(false);
    }
  };

  const handleAddAddress = () => {
    setEditingAddress(undefined);
    setAddressDialogOpen(true);
  };

  const handleEditAddress = (address: UserAddress) => {
    setEditingAddress(address);
    setAddressDialogOpen(true);
  };

  const handleDeleteAddress = async (addressId: string) => {
    if (!confirm("Bạn có chắc chắn muốn xóa địa chỉ này?")) {
      return;
    }

    try {
      const token = getCookieValues<string>(ACCESS_TOKEN);
      const response = await fetch(
        `${API.base_vinh}${API.user.deleteAddress}/${addressId}`,
        {
          method: 'DELETE',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        }
      );

      const data = await response.json();
      
      if (data.code === 10000) {
        toast({
          title: "Thành công",
          description: "Xóa địa chỉ thành công",
        });
        loadAddresses();
      }
    } catch (error: any) {
      console.error("Error deleting address:", error);
      toast({
        title: "Lỗi",
        description: "Không thể xóa địa chỉ",
        variant: "destructive",
      });
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString("vi-VN", {
      day: "2-digit",
      month: "2-digit",
      year: "numeric",
    });
  };

  const handleEditProfile = () => {
    setEditedProfile({
      firstName: profile?.firstName,
      lastName: profile?.lastName,
      dob: profile?.dob,
      phone_number: profile?.phone_number,
      gender: profile?.gender,
    });
    setIsEditing(true);
  };

  const handleCancelEdit = () => {
    setEditedProfile({});
    setIsEditing(false);
  };

  const handleUpdateProfile = async () => {
    if (!profile) return;

    try {
      setIsSaving(true);

      const updateData = {
        firstName: editedProfile.firstName || profile.firstName,
        lastName: editedProfile.lastName || profile.lastName,
        dob: editedProfile.dob || profile.dob,
        phoneNumber: editedProfile.phone_number || profile.phone_number,
        gender: editedProfile.gender === "Nam" ? true : false,
      };

      console.log("🔄 Updating profile with data:", updateData);
      console.log("📍 Profile ID:", profile.id);

      const response = await apiClient.put(
        `${API.user.updateProfile}/${profile.id}`,
        updateData
      );

      console.log("✅ Update response:", response.data);

      if (response.data.code === 10000) {
        toast({
          title: "Cập nhật thành công!",
          description: "Thông tin của bạn đã được cập nhật",
        });

        // Fetch latest profile to update UI and localStorage
        await fetchLatestUserProfile();
        setIsEditing(false);
        setEditedProfile({});
      } else {
        throw new Error(response.data.message || "Cập nhật thất bại");
      }
    } catch (error: any) {
      console.error("❌ Error updating profile:", error);
      toast({
        title: "Lỗi",
        description: error.message || "Không thể cập nhật thông tin",
        variant: "destructive",
      });
    } finally {
      setIsSaving(false);
    }
  };

  if (isLoading) {
    return (
      <div className="container mx-auto px-4 py-16">
        <div className="flex flex-col items-center justify-center gap-4">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
          <p className="text-sm text-muted-foreground">Đang tải thông tin...</p>
        </div>
      </div>
    );
  }

  if (!profile) {
    return (
      <div className="container mx-auto px-4 py-16">
        <div className="text-center">
          <p className="text-lg text-muted-foreground">Không tìm thấy thông tin người dùng</p>
          <Button onClick={() => router.push(ROUTER.home)} className="mt-4">
            Về trang chủ
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="max-w-6xl mx-auto">
        <div className="mb-8">
          <h1 className="text-3xl font-bold mb-2">Thông tin tài khoản</h1>
          <p className="text-muted-foreground">Quản lý thông tin cá nhân và địa chỉ của bạn</p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-1">
            <Card>
              <CardHeader>
                <CardTitle>Ảnh đại diện</CardTitle>
              </CardHeader>
              <CardContent className="flex flex-col items-center">
                <div className="relative mb-4">
                  <img
                    src={profile.image?.valid ? profile.image.data : '/default-avatar.png'}
                    alt="Avatar"
                    className="w-32 h-32 rounded-full object-cover border-4 border-primary/20"
                  />
                </div>
                <h3 className="text-xl font-semibold text-center mb-1">{profile.name}</h3>
                <p className="text-sm text-muted-foreground text-center mb-4">ID: {profile.userId.slice(0, 8)}...</p>
                <Button variant="outline" className="w-full">
                  <Edit className="h-4 w-4 mr-2" />
                  Thay đổi ảnh
                </Button>
              </CardContent>
            </Card>
          </div>

          <div className="lg:col-span-2 space-y-6">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between">
                <CardTitle className="flex items-center gap-2">
                  <User className="h-5 w-5" />
                  Thông tin cá nhân
                </CardTitle>
                {!isEditing ? (
                  <Button variant="ghost" size="sm" onClick={handleEditProfile}>
                    <Edit className="h-4 w-4 mr-2" />
                    Chỉnh sửa
                  </Button>
                ) : (
                  <div className="flex gap-2">
                    <Button variant="outline" size="sm" onClick={handleCancelEdit} disabled={isSaving}>
                      <X className="h-4 w-4 mr-2" />
                      Hủy
                    </Button>
                    <Button variant="default" size="sm" onClick={handleUpdateProfile} disabled={isSaving}>
                      {isSaving ? (
                        <>
                          <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                          Đang lưu...
                        </>
                      ) : (
                        <>
                          <Save className="h-4 w-4 mr-2" />
                          Lưu
                        </>
                      )}
                    </Button>
                  </div>
                )}
              </CardHeader>
              <CardContent className="space-y-4">
                {!isEditing ? (
                  <>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div className="space-y-1">
                        <label className="text-sm font-medium text-muted-foreground">Họ</label>
                        <p className="text-base">{profile.firstName}</p>
                      </div>
                      <div className="space-y-1">
                        <label className="text-sm font-medium text-muted-foreground">Tên</label>
                        <p className="text-base">{profile.lastName}</p>
                      </div>
                    </div>
                    
                    <Separator />
                    
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div className="flex items-center gap-3">
                        <Calendar className="h-5 w-5 text-primary" />
                        <div>
                          <label className="text-sm font-medium text-muted-foreground">Ngày sinh</label>
                          <p className="text-base">{formatDate(profile.dob)}</p>
                        </div>
                      </div>
                      
                      <div className="flex items-center gap-3">
                        <User className="h-5 w-5 text-primary" />
                        <div>
                          <label className="text-sm font-medium text-muted-foreground">Giới tính</label>
                          <p className="text-base">{profile.gender}</p>
                        </div>
                      </div>
                    </div>

                    <Separator />

                    <div className="flex items-center gap-3">
                      <Phone className="h-5 w-5 text-primary" />
                      <div>
                        <label className="text-sm font-medium text-muted-foreground">Số điện thoại</label>
                        <p className="text-base">{profile.phone_number}</p>
                      </div>
                    </div>
                  </>
                ) : (
                  <>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div className="space-y-2">
                        <Label htmlFor="firstName">Họ</Label>
                        <Input
                          id="firstName"
                          value={editedProfile.firstName || ""}
                          onChange={(e) =>
                            setEditedProfile({ ...editedProfile, firstName: e.target.value })
                          }
                          placeholder="Nhập họ"
                          disabled={isSaving}
                        />
                      </div>
                      <div className="space-y-2">
                        <Label htmlFor="lastName">Tên</Label>
                        <Input
                          id="lastName"
                          value={editedProfile.lastName || ""}
                          onChange={(e) =>
                            setEditedProfile({ ...editedProfile, lastName: e.target.value })
                          }
                          placeholder="Nhập tên"
                          disabled={isSaving}
                        />
                      </div>
                    </div>

                    <Separator />

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div className="space-y-2">
                        <Label htmlFor="dob">Ngày sinh</Label>
                        <Input
                          id="dob"
                          type="date"
                          value={editedProfile.dob || ""}
                          onChange={(e) =>
                            setEditedProfile({ ...editedProfile, dob: e.target.value })
                          }
                          disabled={isSaving}
                          max={new Date().toISOString().split('T')[0]}
                        />
                      </div>
                      <div className="space-y-2">
                        <Label htmlFor="gender">Giới tính</Label>
                        <Select
                          value={editedProfile.gender || profile.gender}
                          onValueChange={(value) =>
                            setEditedProfile({ ...editedProfile, gender: value })
                          }
                          disabled={isSaving}
                        >
                          <SelectTrigger id="gender">
                            <SelectValue placeholder="Chọn giới tính" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value="Nam">Nam</SelectItem>
                            <SelectItem value="Nữ">Nữ</SelectItem>
                          </SelectContent>
                        </Select>
                      </div>
                    </div>

                    <Separator />

                    <div className="space-y-2">
                      <Label htmlFor="phone">Số điện thoại</Label>
                      <Input
                        id="phone"
                        type="tel"
                        value={editedProfile.phone_number || ""}
                        onChange={(e) =>
                          setEditedProfile({ ...editedProfile, phone_number: e.target.value })
                        }
                        placeholder="Nhập số điện thoại"
                        disabled={isSaving}
                        onKeyPress={(e) => {
                          if (!/[0-9]/.test(e.key)) {
                            e.preventDefault();
                          }
                        }}
                      />
                    </div>
                  </>
                )}
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between">
                <CardTitle className="flex items-center gap-2">
                  <MapPin className="h-5 w-5" />
                  Địa chỉ ({addresses.length})
                </CardTitle>
                <Button variant="outline" size="sm" onClick={handleAddAddress}>
                  <Plus className="h-4 w-4 mr-2" />
                  Thêm địa chỉ mới
                </Button>
              </CardHeader>
              <CardContent>
                {isLoadingAddresses ? (
                  <div className="flex items-center justify-center py-8">
                    <Loader2 className="h-6 w-6 animate-spin text-primary" />
                  </div>
                ) : addresses && addresses.length > 0 ? (
                  <div className="space-y-4">
                    {addresses.map((address, index) => (
                      <Card key={address.id} className="border-2 hover:border-primary/50 transition-colors">
                        <CardContent className="p-4">
                          <div className="flex items-start justify-between mb-3">
                            <div className="flex items-center gap-2">
                              <h4 className="font-semibold">{address.name}</h4>
                              {index === 0 && (
                                <Badge variant="default">Mặc định</Badge>
                              )}
                            </div>
                            <div className="flex gap-2">
                              <Button 
                                variant="ghost" 
                                size="sm"
                                onClick={() => handleEditAddress(address)}
                              >
                                <Edit className="h-4 w-4" />
                              </Button>
                              <Button 
                                variant="ghost" 
                                size="sm"
                                onClick={() => handleDeleteAddress(address.id)}
                                className="text-destructive hover:text-destructive"
                              >
                                <Trash2 className="h-4 w-4" />
                              </Button>
                            </div>
                          </div>
                          
                          <div className="space-y-2 text-sm">
                            <div className="flex items-center gap-2">
                              <Phone className="h-4 w-4 text-muted-foreground" />
                              <span>{address.phoneNumber}</span>
                            </div>
                            <div className="flex items-start gap-2">
                              <MapPin className="h-4 w-4 text-muted-foreground mt-0.5" />
                              <div className="flex-1">
                                <p className="text-muted-foreground">{address.address.other}</p>
                                <p className="text-xs text-muted-foreground mt-1">
                                  {address.address.ward.fullName}, {address.address.ward.district.fullName}, {address.address.ward.district.province.fullName}
                                </p>
                              </div>
                            </div>
                          </div>
                        </CardContent>
                      </Card>
                    ))}
                  </div>
                ) : (
                  <div className="text-center py-8 text-muted-foreground">
                    <MapPin className="h-12 w-12 mx-auto mb-3 opacity-20" />
                    <p>Chưa có địa chỉ nào</p>
                    <Button variant="outline" className="mt-4" onClick={handleAddAddress}>
                      <Plus className="h-4 w-4 mr-2" />
                      Thêm địa chỉ đầu tiên
                    </Button>
                  </div>
                )}
              </CardContent>
            </Card>
          </div>
        </div>

        <AddressDialog
          open={addressDialogOpen}
          onOpenChange={setAddressDialogOpen}
          onSuccess={loadAddresses}
          editAddress={editingAddress}
        />
      </div>
    </div>
  );
}
