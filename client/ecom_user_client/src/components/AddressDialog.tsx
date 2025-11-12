"use client";

import { useState, useEffect } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Loader2, Check, ChevronsUpDown } from "lucide-react";
import { cn } from "@/lib/utils";
import API from "@/assets/configs/api";
import { useToast } from "@/hooks/use-toast";
import type { Province, District, Ward, AddressFormData, UserAddress } from "@/types/address.types";
import { cookies } from "@/assets/helpers";
import { ACCESS_TOKEN } from "@/assets/configs/request";

interface AddressDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess: () => void;
  editAddress?: UserAddress;
}

export default function AddressDialog({
  open,
  onOpenChange,
  onSuccess,
  editAddress,
}: AddressDialogProps) {
  const { toast } = useToast();
  const [isLoading, setIsLoading] = useState(false);

  // Form state
  const [name, setName] = useState("");
  const [phoneNumber, setPhoneNumber] = useState("");
  const [detailAddress, setDetailAddress] = useState("");

  // Address selection state
  const [provinces, setProvinces] = useState<Province[]>([]);
  const [districts, setDistricts] = useState<District[]>([]);
  const [wards, setWards] = useState<Ward[]>([]);

  const [selectedProvince, setSelectedProvince] = useState<Province | null>(null);
  const [selectedDistrict, setSelectedDistrict] = useState<District | null>(null);
  const [selectedWard, setSelectedWard] = useState<Ward | null>(null);

  // Popover state
  const [openProvince, setOpenProvince] = useState(false);
  const [openDistrict, setOpenDistrict] = useState(false);
  const [openWard, setOpenWard] = useState(false);

  // Load provinces on mount
  useEffect(() => {
    if (open) {
      loadProvinces();
    }
  }, [open]);

  // Load edit data
  useEffect(() => {
    if (editAddress && open) {
      setName(editAddress.name);
      setPhoneNumber(editAddress.phoneNumber);
      setDetailAddress(editAddress.address.other);
      
      // Set selected province, district, ward
      const ward = editAddress.address.ward;
      const district = ward.district;
      const province = district.province;
      
      setSelectedProvince(province);
      setSelectedDistrict(district);
      setSelectedWard(ward);
      
      // Load districts and wards
      loadDistricts(province.id);
      loadWards(district.id);
    } else if (!editAddress && open) {
      // Reset form
      resetForm();
    }
  }, [editAddress, open]);

  const resetForm = () => {
    setName("");
    setPhoneNumber("");
    setDetailAddress("");
    setSelectedProvince(null);
    setSelectedDistrict(null);
    setSelectedWard(null);
    setDistricts([]);
    setWards([]);
  };

  const loadProvinces = async () => {
    try {
      const response = await fetch(API.base_vinh + API.address.provinces, {
        headers: {
          'Content-Type': 'application/json',
        },
      });
      const data = await response.json();
      if (data.code === 10000) {
        setProvinces(data.result);
      }
    } catch (error) {
      console.error("Error loading provinces:", error);
      toast({
        title: "Lỗi",
        description: "Không thể tải danh sách tỉnh/thành phố",
        variant: "destructive",
      });
    }
  };

  const loadDistricts = async (provinceId: string) => {
    try {
      const response = await fetch(
        `${API.base_vinh}${API.address.districts}?province-id=${provinceId}`,
        {
          headers: {
            'Content-Type': 'application/json',
          },
        }
      );
      const data = await response.json();
      if (data.code === 10000) {
        setDistricts(data.result);
      }
    } catch (error) {
      console.error("Error loading districts:", error);
      toast({
        title: "Lỗi",
        description: "Không thể tải danh sách quận/huyện",
        variant: "destructive",
      });
    }
  };

  const loadWards = async (districtId: string) => {
    try {
      const response = await fetch(
        `${API.base_vinh}${API.address.wards}?district-id=${districtId}`,
        {
          headers: {
            'Content-Type': 'application/json',
          },
        }
      );
      const data = await response.json();
      if (data.code === 10000) {
        setWards(data.result);
      }
    } catch (error) {
      console.error("Error loading wards:", error);
      toast({
        title: "Lỗi",
        description: "Không thể tải danh sách phường/xã",
        variant: "destructive",
      });
    }
  };

  const handleProvinceSelect = (province: Province) => {
    setSelectedProvince(province);
    setSelectedDistrict(null);
    setSelectedWard(null);
    setDistricts([]);
    setWards([]);
    loadDistricts(province.id);
    setOpenProvince(false);
  };

  const handleDistrictSelect = (district: District) => {
    setSelectedDistrict(district);
    setSelectedWard(null);
    setWards([]);
    loadWards(district.id);
    setOpenDistrict(false);
  };

  const handleWardSelect = (ward: Ward) => {
    setSelectedWard(ward);
    setOpenWard(false);
  };

  const handleSubmit = async () => {
    // Validation
    if (!name.trim()) {
      toast({
        title: "Lỗi",
        description: "Vui lòng nhập tên người nhận",
        variant: "destructive",
      });
      return;
    }

    if (!phoneNumber.trim() || !/^[0-9]{10,11}$/.test(phoneNumber)) {
      toast({
        title: "Lỗi",
        description: "Vui lòng nhập số điện thoại hợp lệ (10-11 số)",
        variant: "destructive",
      });
      return;
    }

    if (!selectedWard) {
      toast({
        title: "Lỗi",
        description: "Vui lòng chọn đầy đủ địa chỉ (Tỉnh/Quận/Xã)",
        variant: "destructive",
      });
      return;
    }

    if (!detailAddress.trim()) {
      toast({
        title: "Lỗi",
        description: "Vui lòng nhập địa chỉ chi tiết",
        variant: "destructive",
      });
      return;
    }

    setIsLoading(true);

    try {
      const payload: AddressFormData = {
        name: name.trim(),
        phoneNumber: phoneNumber.trim(),
        address: {
          wardId: selectedWard.id,
          other: detailAddress.trim(),
        },
      };

      const url = editAddress
        ? `${API.base_vinh}${API.user.updateAddress}/${editAddress.id}`
        : `${API.base_vinh}${API.user.createAddress}`;

      const method = editAddress ? 'PUT' : 'POST';

      // Get token from cookies
      const token = cookies.getCookieValues(ACCESS_TOKEN)
      // console.log("Submitting address with payload:", token);
      const response = await fetch(url, {
        method,
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(payload),
      });

      const data = await response.json();

      if (data.code === 10000) {
        toast({
          title: "Thành công",
          description: editAddress
            ? "Cập nhật địa chỉ thành công"
            : "Thêm địa chỉ mới thành công",
        });
        onSuccess();
        onOpenChange(false);
        resetForm();
      } else {
        throw new Error(data.message || "API returned error");
      }
    } catch (error: any) {
      console.error("Error submitting address:", error);
      toast({
        title: "Lỗi",
        description: error.message || "Không thể lưu địa chỉ",
        variant: "destructive",
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>
            {editAddress ? "Cập nhật địa chỉ" : "Thêm địa chỉ mới"}
          </DialogTitle>
        </DialogHeader>

        <div className="space-y-4 py-4">
          {/* Name */}
          <div className="space-y-2">
            <Label htmlFor="name">
              Tên người nhận <span className="text-destructive">*</span>
            </Label>
            <Input
              id="name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Nhập tên người nhận"
              disabled={isLoading}
            />
          </div>

          {/* Phone Number */}
          <div className="space-y-2">
            <Label htmlFor="phone">
              Số điện thoại <span className="text-destructive">*</span>
            </Label>
            <Input
              id="phone"
              type="tel"
              value={phoneNumber}
              onChange={(e) => setPhoneNumber(e.target.value)}
              placeholder="Nhập số điện thoại"
              disabled={isLoading}
            />
          </div>

          {/* Province */}
          <div className="space-y-2">
            <Label>
              Tỉnh/Thành phố <span className="text-destructive">*</span>
            </Label>
            <Popover open={openProvince} onOpenChange={setOpenProvince}>
              <PopoverTrigger asChild>
                <Button
                  variant="outline"
                  role="combobox"
                  aria-expanded={openProvince}
                  className="w-full justify-between"
                  disabled={isLoading}
                >
                  {selectedProvince
                    ? selectedProvince.fullName
                    : "Chọn tỉnh/thành phố"}
                  <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-full p-0" align="start">
                <Command>
                  <CommandInput placeholder="Tìm kiếm tỉnh/thành phố..." />
                  <CommandList>
                    <CommandEmpty>Không tìm thấy kết quả.</CommandEmpty>
                    <CommandGroup>
                      {provinces.map((province) => (
                        <CommandItem
                          key={province.id}
                          value={province.fullName}
                          onSelect={() => handleProvinceSelect(province)}
                        >
                          <Check
                            className={cn(
                              "mr-2 h-4 w-4",
                              selectedProvince?.id === province.id
                                ? "opacity-100"
                                : "opacity-0"
                            )}
                          />
                          {province.fullName}
                        </CommandItem>
                      ))}
                    </CommandGroup>
                  </CommandList>
                </Command>
              </PopoverContent>
            </Popover>
          </div>

          {/* District */}
          <div className="space-y-2">
            <Label>
              Quận/Huyện <span className="text-destructive">*</span>
            </Label>
            <Popover open={openDistrict} onOpenChange={setOpenDistrict}>
              <PopoverTrigger asChild>
                <Button
                  variant="outline"
                  role="combobox"
                  aria-expanded={openDistrict}
                  className="w-full justify-between"
                  disabled={!selectedProvince || isLoading}
                >
                  {selectedDistrict
                    ? selectedDistrict.fullName
                    : "Chọn quận/huyện"}
                  <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-full p-0" align="start">
                <Command>
                  <CommandInput placeholder="Tìm kiếm quận/huyện..." />
                  <CommandList>
                    <CommandEmpty>Không tìm thấy kết quả.</CommandEmpty>
                    <CommandGroup>
                      {districts.map((district) => (
                        <CommandItem
                          key={district.id}
                          value={district.fullName}
                          onSelect={() => handleDistrictSelect(district)}
                        >
                          <Check
                            className={cn(
                              "mr-2 h-4 w-4",
                              selectedDistrict?.id === district.id
                                ? "opacity-100"
                                : "opacity-0"
                            )}
                          />
                          {district.fullName}
                        </CommandItem>
                      ))}
                    </CommandGroup>
                  </CommandList>
                </Command>
              </PopoverContent>
            </Popover>
          </div>

          {/* Ward */}
          <div className="space-y-2">
            <Label>
              Phường/Xã <span className="text-destructive">*</span>
            </Label>
            <Popover open={openWard} onOpenChange={setOpenWard}>
              <PopoverTrigger asChild>
                <Button
                  variant="outline"
                  role="combobox"
                  aria-expanded={openWard}
                  className="w-full justify-between"
                  disabled={!selectedDistrict || isLoading}
                >
                  {selectedWard ? selectedWard.fullName : "Chọn phường/xã"}
                  <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                </Button>
              </PopoverTrigger>
              <PopoverContent className="w-full p-0" align="start">
                <Command>
                  <CommandInput placeholder="Tìm kiếm phường/xã..." />
                  <CommandList>
                    <CommandEmpty>Không tìm thấy kết quả.</CommandEmpty>
                    <CommandGroup>
                      {wards.map((ward) => (
                        <CommandItem
                          key={ward.id}
                          value={ward.fullName}
                          onSelect={() => handleWardSelect(ward)}
                        >
                          <Check
                            className={cn(
                              "mr-2 h-4 w-4",
                              selectedWard?.id === ward.id
                                ? "opacity-100"
                                : "opacity-0"
                            )}
                          />
                          {ward.fullName}
                        </CommandItem>
                      ))}
                    </CommandGroup>
                  </CommandList>
                </Command>
              </PopoverContent>
            </Popover>
          </div>

          {/* Detail Address */}
          <div className="space-y-2">
            <Label htmlFor="detail">
              Địa chỉ chi tiết <span className="text-destructive">*</span>
            </Label>
            <Input
              id="detail"
              value={detailAddress}
              onChange={(e) => setDetailAddress(e.target.value)}
              placeholder="Số nhà, tên đường, ấp/thôn..."
              disabled={isLoading}
            />
          </div>
        </div>

        <DialogFooter>
          <Button
            variant="outline"
            onClick={() => onOpenChange(false)}
            disabled={isLoading}
          >
            Hủy
          </Button>
          <Button onClick={handleSubmit} disabled={isLoading}>
            {isLoading ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Đang lưu...
              </>
            ) : (
              editAddress ? "Cập nhật" : "Thêm mới"
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
