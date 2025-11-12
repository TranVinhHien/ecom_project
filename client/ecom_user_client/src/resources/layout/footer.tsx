import { useTranslations } from "next-intl";
import React from "react";

export default function Footer() {
  const t = useTranslations("footer")
  return (
    <footer className="bg-background pt-8 text-foreground">
      {/* Top section */}
      <div className="max-w-7xl mx-auto px-6 grid grid-cols-1 md:grid-cols-4 gap-8 pb-8 border-b border-[hsl(var(--primary)/0.15)]">
        {/* VỀ CHÚNG TÔI */}
        <div>
          <div className="font-bold mb-3 text-[hsl(var(--primary))]">{t("ve_chung_toi")}</div>
          <ul className="space-y-2 text-muted-foreground text-sm">
            <li><a href="#">{t("gioi_thieu_ecom_vn")}</a></li>
            <li><a href="#">{t("gioi_thieu_ecommal")}</a></li>
            <li><a href="#">{t("quy_che_hoat_dong")}</a></li>
            <li><a href="#">{t("chinh_sach_bao_mat")}</a></li>
            <li><a href="#">{t("giao_hang_va_nhan_hang")}</a></li>
            <li><a href="#">{t("dieu_khoan_su_dung")}</a></li>
          </ul>
        </div>
        {/* DÀNH CHO NGƯỜI MUA */}
        <div>
          <div className="font-bold mb-3 text-[hsl(var(--primary))]">{t("danh_cho_nguoi_mua")}</div>
          <ul className="space-y-2 text-muted-foreground text-sm">
            <li><a href="#">{t("giai_quyet_kieu_nai")}</a></li>
            <li><a href="#">{t("huong_dan_mua_hang")}</a></li>
            <li><a href="#">{t("chinh_sach_doi_tra")}</a></li>
            <li><a href="#">{t("cham_soc_khach_hang")}</a></li>
            <li><a href="#">{t("nap_tien_dien_thoai")}</a></li>
          </ul>
        </div>
        {/* DÀNH CHO NGƯỜI BÁN */}
        <div>
          <div className="font-bold mb-3 text-[hsl(var(--primary))]">{t("danh_cho_don_vi_van_chuyen")}</div>
          <ul className="space-y-2 text-muted-foreground text-sm">
            <li><a href="#">{t("quy_dinh_doi_voi_don_vi_van_chuyen")}</a></li>
            <li><a href="#">{t("chinh_sach_ban_hang")}</a></li>
            <li><a href="#">{t("he_thong_tieu_chi_kiem_duyet")}</a></li>
          </ul>
        </div>
        {/* TẢI ỨNG DỤNG */}
        <div>
          <div className="font-bold mb-3 text-[hsl(var(--primary))]">{t("tai_ung_dung_ecom")}</div>
          <div className="text-sm text-muted-foreground mb-2">{t("mang_the_gioi_mua_sam_cua_ecom")}<br />{t("trong_tam_tay_ban")}</div>
          <div className="flex flex-col gap-2">
            <div className="flex gap-2">
              <img src="/appstore.png" alt="App Store" className="h-10" />
              <img src="/googleplay.png" alt="Google Play" className="h-10" />
            </div>
            <img src="/appgallery.png" alt="App Gallery" className="h-10 w-60" />
          </div>
        </div>
      </div>

      {/* Middle section */}
      <div className="bg-[hsl(var(--primary))] text-[hsl(var(--primary-foreground))] py-8">
        <div className="max-w-7xl mx-auto px-6 grid grid-cols-1 md:grid-cols-2 gap-8">
          {/* Company info */}
          <div>
            <div className="font-bold mb-2">Công ty trách nhiệm cổ phần Hienlazada tập đoàn HUIT</div>
            <div className="text-sm mb-2 text-[hsl(var(--primary-foreground))]">
             {t("so_dkkd")}: {t("so_dkkd_value")}<br />
              {t("co_quan_cap")}: {t("co_quan_cap_value")}
            </div>
            <div className="text-sm mb-2 text-[hsl(var(--primary-foreground))]">
              {t("dia_chi")}: {t("dia_chi_value")}
            </div>
            <div className="text-sm mb-2 text-[hsl(var(--primary-foreground))]">{t("email")}: vinhhien12z@gmail.com</div>
            <div className="flex gap-4 mt-4">
              <img src="/bocongthuong1.png" alt="Bộ Công Thương" className="h-10" />
              <img src="/bocongthuong2.png" alt="Bộ Công Thương" className="h-10" />
            </div>
          </div>
          {/* Đăng ký nhận tin */}
          <div className="flex flex-col items-start justify-center">
            <div className="font-bold mb-2">{t("dang_ky_nhan_tin")}</div>
            <form className="flex w-full max-w-md">
              <input
                type="email"
                placeholder="Email của bạn là"
                className="flex-1 px-4 py-2 rounded-l bg-white text-gray-800 outline-none"
              />
              <button
                type="submit"
                className="px-6 py-2 bg-[#ee4d2d] text-white font-bold rounded-r hover:bg-[#d84315] transition"
              >
                {t("dang_ky")}
              </button>
            </form>
          </div>
        </div>
      </div>
    </footer>
  );
}
