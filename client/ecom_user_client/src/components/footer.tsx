import React from "react";

export default function Footer() {
  return (
    <footer className="bg-[#f7f7f7] pt-8">
      {/* Top section */}
      <div className="max-w-7xl mx-auto px-6 grid grid-cols-1 md:grid-cols-4 gap-8 pb-8 border-b border-gray-200">
        {/* Về chúng tôi */}
        <div>
          <div className="font-bold mb-3 text-gray-800">VỀ CHÚNG TÔI</div>
          <ul className="space-y-2 text-gray-700 text-sm">
            <li><a href="#">Giới thiệu Sendo.vn</a></li>
            <li><a href="#">Giới thiệu SenMall</a></li>
            <li><a href="#">Quy chế hoạt động</a></li>
            <li><a href="#">Chính sách bảo mật</a></li>
            <li><a href="#">Giao hàng và Nhận hàng</a></li>
            <li><a href="#">Điều khoản sử dụng</a></li>
          </ul>
        </div>
        {/* Dành cho người mua */}
        <div>
          <div className="font-bold mb-3 text-gray-800">DÀNH CHO NGƯỜI MUA</div>
          <ul className="space-y-2 text-gray-700 text-sm">
            <li><a href="#">Giải quyết khiếu nại</a></li>
            <li><a href="#">Hướng dẫn mua hàng</a></li>
            <li><a href="#">Chính sách đổi trả</a></li>
            <li><a href="#">Chăm sóc khách hàng</a></li>
            <li><a href="#">Nạp tiền điện thoại</a></li>
          </ul>
        </div>
        {/* Dành cho người bán */}
        <div>
          <div className="font-bold mb-3 text-gray-800">DÀNH CHO NGƯỜI BÁN</div>
          <ul className="space-y-2 text-gray-700 text-sm">
            <li><a href="#">Quy định đối với người bán</a></li>
            <li><a href="#">Chính sách bán hàng</a></li>
            <li><a href="#">Hệ thống tiêu chí kiểm duyệt</a></li>
            <li><a href="#">Mở shop trên Sendo</a></li>
          </ul>
        </div>
        {/* Tải ứng dụng */}
        <div>
          <div className="font-bold mb-3 text-gray-800">TẢI ỨNG DỤNG SENDO</div>
          <div className="text-sm text-gray-700 mb-2">Mang thế giới mua sắm của Sendo<br />trong tầm tay bạn</div>
          <div className="flex flex-col gap-2">
            <div className="flex gap-2">
              
              <img src="/appstore.png" alt="App Store" className="h-10" />
              <img src="/googleplay.png" alt="Google Play" className="h-10" />
            </div>
          </div>
        </div>
      </div>

      {/* Middle section */}
      <div className="bg-[#222f3e] text-white py-8">
        <div className="max-w-7xl mx-auto px-6 grid grid-cols-1 md:grid-cols-2 gap-8">
          {/* Company info */}
          <div>
            <div className="font-bold mb-2">Công ty Cổ phần Công nghệ Sen Đỏ, thành viên của Tập đoàn FPT</div>
            <div className="text-sm mb-2">
              Số ĐKKD: 0312776486 - Ngày cấp: 13/05/2014, được sửa đổi lần thứ 20, ngày 26/04/2022.<br />
              Cơ quan cấp: Sở Kế hoạch và Đầu tư TPHCM.
            </div>
            <div className="text-sm mb-2">
              Địa chỉ: Tầng 5, Tòa nhà A, Vườn Ươm Doanh Nghiệp, Lô D.01, Đường Tân Thuận, Khu chế xuất Tân Thuận, Phường Tân Thuận Đông, Quận 7, Thành phố Hồ Chí Minh, Việt Nam.
            </div>
            <div className="text-sm mb-2">Email: lienhe@sendo.vn</div>
            <div className="flex gap-4 mt-4">
              <img src="/bocongthuong1.png" alt="Bộ Công Thương" className="h-10" />
            </div>
          </div>
          {/* Đăng ký nhận tin */}
          <div className="flex flex-col items-start justify-center">
            <div className="font-bold mb-2">Đăng ký nhận bản tin ưu đãi khủng từ Sendo</div>
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
                Đăng ký
              </button>
            </form>
          </div>
        </div>
      </div>
    </footer>
  );
} 