import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { randomItem } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

// ==============================================================================
// 1. CẤU HÌNH & DỮ LIỆU ĐẦU VÀO
// ==============================================================================
import crypto from 'k6/crypto'; // ĐÚNG
import { b64encode } from 'k6/encoding';
// Danh sách Category Path (Để random)
const CATE_PATHS = [
    '/add-girl-discover-290',
    '/add-girl-discover-290/manage-near-final-4392',
    '/add-girl-discover-290/trial-stay-case-2878',
    '/change-four-view-551',
    '/change-four-view-551/decade-structure-6299',
    '/live-full-wish-into-732/defense-level-1594',
    '/live-full-wish-into-732/hot-lead-still-hear-8345',
    '/live-full-wish-into-732/join-spring-soon-4608',
    '/technology-the-343',
];

const SECRET_KEY = 'bv-T"-u6@-WR?SHiHQ7yQ]CK*dd9(@jM9BI)|g;zq)ur-Z.Jw/u5HyJHgg,KS.fa';

// --- CÁC HÀM TIỆN ÍCH (HELPER) ---

// 1. Hàm chuyển đổi Base64 thường sang Base64URL (Chuẩn JWT bắt buộc)
function toBase64Url(base64String) {
    return base64String.replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
}

// 2. Hàm tạo UUID v4 (Giả lập)
function uuidv4() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        var r = Math.random() * 16 | 0, v = c == 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
    });
}

// 3. Hàm tạo chuỗi ngẫu nhiên (cho username)
function randomString(length) {
    const chars = 'abcdefghijklmnopqrstuvwxyz0123456789';
    let result = '';
    for (let i = 0; i < length; i++) {
        result += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    return result;
}

// --- HÀM TẠO TOKEN CHÍNH ---
export  function generateJWT() {
    // 1. Header (HS512)
    const header = {
        alg: "HS512",
        typ: "JWT"
    };
    const headerBase64 = toBase64Url(b64encode(JSON.stringify(header)));

    // 2. Payload (Dữ liệu động)
    const now = Math.floor(Date.now() / 1000);
    const exp = now + 3600;

    const payload = {
        "sub": randomString(8),
        "scope": "ROLE_USER",
        "iss": "lemarchenoble.id.vn",
        "exp": exp,
        "iat": now,
        "userId": uuidv4(),
        "jti": uuidv4(),
        "email": `${randomString(10)}@gmail.com`
    };
    const payloadBase64 = toBase64Url(b64encode(JSON.stringify(payload)));

    // 3. Chữ ký (Signature) - Dùng HS512
    const signatureInput = `${headerBase64}.${payloadBase64}`;
    
    // --- SỬA LẠI CÁCH GỌI HMAC ---
    // Dùng trực tiếp biến 'crypto' đã import
    const signature = crypto.hmac('sha512', SECRET_KEY, signatureInput, 'base64');
    
    const signatureBase64Url = toBase64Url(signature);

    // 4. Ghép lại thành Token hoàn chỉnh
    return `${headerBase64}.${payloadBase64}.${signatureBase64Url}`;
}


// Danh sách Shop ID (Để random)
const SHOP_IDS = [
    'shop_1',
    'shop_2',
    'shop_3',
    'shop_6',
    'shop_5',
    'shop_7',
    'shop_8',
    'shop_9',
    'shop_0',
    'shop_1',
    'shop_12',
    'shop_13',
    'shop_14',
    'shop_15',
    'shop_16',
    'shop_17',
    'shop_18',
    'shop_19',
    'shop_20',
    'shop_21',
];

// Danh sách Sort Options
const SORT_OPTIONS = ["price_asc", "price_desc", "name_asc", "name_desc", "best_sell"];

// Hàm random số nguyên trong khoảng [min, max]
function randomInt(min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}


// Mảng Product Key (Dùng cho API Get Detail)
const PRODUCT_KEYS =["happen-indeed-do-fe7b8937", "pick-success-change-44ea8a66", "without-or-serve-e20d3119", "account-imagine-05596135", "vote-can-news-1c8c634f", "network-pm-hit-4ab24388", "add-commercial-sort-3ce60f05", "ground-number-less-00f3db36", "best-area-strategy-a7b228db", "pattern-class-on-78f2e27e", "understand-not-f43e6c66", "outside-republican-91257b10", "meeting-avoid-event-554a06b4", "there-behind-sister-f96cb463", "occur-program-use-e0a129a6", "education-report-30ed4acc", "stop-appear-453bfb83", "week-claim-term-cab62c38", "player-customer-7d303978", "final-away-easy-b81c573c", "brother-i-director-93e58af0", "close-business-6294e1ef", "itself-send-4501f62a", "another-language-1e848181", "consumer-direction-27918abc", "experience-project-bdc244b7", "behind-likely-oil-44073f24", "scene-real-outside-4d1c784c", "never-third-my-bag-a128de0c", "political-98a2f8e1", "drop-treat-decide-d79fed9b", "make-born-trouble-ccf4c0f5", "partner-later-a7fd21f2", "on-in-plant-world-511bd9d7", "west-war-another-0fa4acf8", "article-item-guy-b302aa17", "indeed-win-key-4f6c889f", "between-carry-dark-ec853fde", "claim-early-cc15c4eb", "official-threat-32357ba9", "yourself-defense-if-59a25d66", "simply-beat-any-16fc4ddb", "during-company-8c4fec8d", "ten-daughter-board-73c5a630", "reveal-guy-ea71ffbf", "hope-pretty-right-76f779c9", "loss-operation-in-a502e9fe", "value-she-half-c2e372f0", "eight-amount-yeah-4cefbd3a", "growth-appear-cut-569224e9", "develop-purpose-06e0dabb", "plant-economy-hope-2bcd8ebf", "relationship-7a6cbdbf", "soon-people-point-3d3f901c", "impact-brother-59542a0f", "low-available-75d5db96", "their-now-heart-for-26dca54b", "choice-production-b452ebc6", "character-personal-830c7019", "future-film-122d30eb", "collection-training-eb717245", "mrs-society-7db9fb89", "it-will-member-ball-73ad280d", "about-walk-box-edf76307", "wonder-back-show-773797d4", "tonight-seem-18d13682", "chair-small-be-11e295cd", "bar-reveal-wonder-ebad2ebd", "class-western-95fffe5b", "beyond-most-area-a08433b5", "major-research-2d1762e3", "candidate-close-1d8c8dc4", "necessary-real-9734a12c", "become-walk-black-f39ffe30", "suggest-right-bad8e7b2", "low-age-rock-here-2be37428", "nothing-hold-family-9d1b332b", "responsibility-deep-26079b53", "threat-feeling-c393b6ea", "another-friend-c2f8749b", "section-prove-would-94cd0c47", "class-degree-0d3481f9", "second-lose-explain-f444cc45", "everyone-movie-f13492ac", "age-feel-husband-49c937fa", "we-free-toward-5bf956b1", "possible-certain-b9598869", "forward-television-c83b5cd4", "free-whole-article-3616c731", "prevent-end-5cbedc67", "near-growth-share-e7612686", "station-person-2f68b86f", "pretty-sport-1ba8ccef", "cut-yet-guess-c36e40cd", "maintain-director-ce61d6ae", "reduce-stuff-push-f9bc1203", "security-degree-2db441d8", "man-young-1e5b2741", "wish-wind-myself-39e47bb0", "take-hot-man-mr-50544039"];

// Mảng SKU ID (Dùng cho API Create Order)
const SKU_IDS = ["18bdd27a-1d4b-4d86-90ca-45afd27dbbfc", "a0dfbae1-7619-4e45-89fe-980ab33aedaa", "82d05bb9-529e-4631-9cf3-bd596ca18631", "356f9bc8-5065-4ef7-b87d-3d4beb09bcf6", "68da2a86-23b4-4eeb-a2ae-d550a12c1569", "8e532f19-50b2-44de-b9ab-ffc304ddad04", "a91a622d-bb1f-4ea9-aa65-d85fd84bf9fa", "269b6ced-11c3-4058-b247-0404f264d613", "b2cc191d-6eb6-4625-a28d-ba6afea5a8c4", "53c267dd-ed06-4353-a7c7-0727e5f28812", "0532da4b-0257-4be9-b1ba-78d5bf4cf0f4", "40c8cd61-eb31-4a31-98b4-7a0762c05497", "4bc2af35-3af6-43ec-bad0-d57953bc7b63", "7f794095-5b7d-4d39-8d18-7bcdf5034f78", "2b74fcfc-6b8b-4c82-852f-c5945a126321", "c45d13c4-eddd-489c-8db2-3672e643773c", "f763ef27-a409-4d2f-a0b2-a5271aacee6d", "b7bab6ae-1c31-466f-b2a1-32bf523880d7", "a5bfd6dc-e823-436c-a438-067ccbd141f4", "5e9a6bdf-9d97-4bc8-a2d7-6ab5e46e0803", "42a1f2f3-a1ad-4339-949c-e32f4fe79f4f", "d2b1e2b8-aa38-488a-8e56-2083bbb329f4", "6721b3c5-e3d6-4e6c-a06f-6840441ad0b4", "b0e7b879-de8d-4085-bc9a-837a4bb8b282", "fa52a3a5-8980-4500-92b1-29c5fdcc9147", "08f624b3-706d-464e-abdd-f53a57fbb88b", "19f65421-1d2c-4395-8d1f-2a972ad4d41b", "93cba776-f35b-4a7c-8e07-8b3a55eb670d", "2bdd772d-0535-4986-94f7-d2ceca617277", "075d8bff-6ed1-412c-aef3-10b19d1be0a7", "99af84b0-0547-446c-ac1f-a336e5da36c7", "7f246ef4-19c2-49a1-8cac-338e0ad05e5b", "027a9c7d-4a3c-48c9-b5d5-29c5a26aad86", "c051db2c-6222-4369-a7c1-0d337a60f931", "65a39f81-3fa8-46c4-aa52-119d1a2ee631", "2b5fb928-ae44-48c5-b651-28be0a2a2cce", "260e3bb4-93a9-4a3e-a28d-4ab87c5ae142", "4f745e47-c8cb-4dfe-8e26-1cede1536366", "135469eb-06fd-4c46-a5c0-4f8105304abd", "1079209f-f8ec-4d39-9ec6-2dd95dd4200f", "a62d2ce5-1791-4ac7-b248-5e39595548ca", "5bb0162e-d036-4971-9c42-2bd88bf01c2c", "b7e6a063-a8d2-4675-a50a-04091e4c7084", "151179a0-9a43-48cd-a66c-eb7a800ccac1", "66a84938-97c7-4a19-8001-3069200c771e", "3d60f0fe-a68c-47ba-bc20-f314da5c739f", "b7e964c4-fad0-443b-9619-7ad7039e9696", "b4f57e99-02a7-47db-aeac-fb158960b71b", "9483fe9c-675d-4a6a-9830-7064a6e2d409", "d7c324f7-d9d4-48fe-892d-bf00c80525d2", "b8db5ee5-45a3-4407-a5be-c438964cd5e8", "e28a586f-aa99-404d-81c4-2cba98a51256", "1894e28c-d71e-4620-b7ce-7f34acd33177", "0a828360-7d3a-486d-be49-b7a6b0bfb58e", "ad66abf2-f84e-4e79-b187-a8b9bb5a68ec", "cd7a9bc5-28c8-44d6-be69-926697767ddf", "f9890f0d-7f8d-4072-81b1-68794db5262a", "76212554-ff5a-4dc5-bf76-f2fe0c704e22", "bf08fc8d-1515-492d-988c-705f7c0407d3", "29c665c0-e83c-4aaf-8810-029806cb79bc", "b6e080e0-75ca-45f4-bdbe-2ddbe98a7446", "1b94bcd4-a481-431d-9c15-c4792155772d", "c369efeb-79d0-4f9e-9b51-f2eac0a01444", "d20a6fdf-3910-4c94-972c-cfb6bc069583", "3ad54cef-86a1-46f3-8b44-13350881f08f", "2c01b109-15ee-4c02-a6aa-1ddf5b6da954", "c67b4bf8-41e1-4398-825b-b41aeb1609d1", "5a706559-f3b4-4e36-be05-02bf2a7a7d21", "c9a153c9-1e71-4916-9537-bd01a20f8f04", "2fc62477-32d3-48d8-93a9-66358334393f", "e35c1ec8-c269-4c50-afed-d6346a0b5c68", "c1e53902-dcaf-4950-9a69-e53dfdfa0bbe", "c668461e-a5c6-44b5-8915-f0bfefe2b753", "2e192bf4-0612-4bf4-8023-eba01b5b25f0", "cc027a83-7d56-4aa1-83a1-3f9bd2626aa6", "01e7afba-b5b6-434f-9302-e19090d38719", "03202f00-4ea5-4bef-b185-80389d719249", "dc1bb87f-879b-4035-b5ba-7f4788c0076f", "3a8d9323-7224-464c-bcbd-21c8f04d8fa4", "2c74fa5e-9fe7-4e89-96ac-c43277800f37", "39c16166-fb95-4b10-a340-28bee6326601", "eb3ebe74-fe7d-454e-a196-972e6e27e227", "16ff2212-0f5b-487d-bd1f-92a1bcf93047", "2319d149-5dcd-4aef-b70d-bbbad50f31b8", "07c55fed-549d-40f4-9732-3a11da63f134", "7e4471cd-e304-462f-b17f-58a43c9675fd", "de4a6401-63ed-4b45-b1a4-8a865e73c629", "2efdfc87-530e-44d9-bf05-694b65f3ee5e", "e1bdc970-c2eb-4137-b241-0a0ee6d2b619", "6aca14f5-377b-4708-99ff-261dee4d0069", "c8b58074-0194-44cc-b120-373e26afe5a5", "330d8f26-6286-4cb6-9b75-bc86d6378830", "d2d4d5cf-0c73-47ea-8c1e-a3fce254ea0d", "08bae50b-82c6-463d-87d6-aa410ad7a8ec", "8c712ca6-c648-4b05-a718-f96f91c55d8a", "bda7022d-f8f6-4434-aa4b-b573354a1448", "c2fce886-ade3-4243-9506-9525afdfbfaa", "6329c011-68c6-49a3-972b-4da7d8ada278", "d54a2fdd-4c83-4b5d-b477-90d5cd3f9097", "d099806d-b616-4cb7-a25e-3ef71fa5600c"];

// Cấu hình URL Service
const BASE_URL_PRODUCT = 'http://172.26.127.95:9001/v1'; 
const BASE_URL_ORDER = 'http://172.26.127.95:9002/v1';

export const options = {
    stages: [
        { duration: '1m', target: 100 },
        { duration: '1m', target: 75 },
        { duration: '30s', target: 25 },
    ],
    thresholds: {
        http_req_duration: ['p(95)<1000'],
        http_req_failed: ['rate<0.01'],
    },
};

// ==============================================================================
// 2. LOGIC TEST
// ==============================================================================

export default function () {
    const rand = Math.random();

    // --- SCENARIO A: GET ALL PRODUCTS (LIST) - 40% ---
    if (rand < 0.5) {
        group('API_Get_All_Product', function () {
            
            // 1. Khởi tạo params cơ bản
            let params = {
                page: 1,
                limit: 21, // Random limit từ 10-30 nếu muốn: Math.floor(Math.random() * 20) + 10
            };

            // 2. Random Category (Luôn chọn 1 cái từ danh sách)
            if (CATE_PATHS.length > 0) {
                params.cate_path = randomItem(CATE_PATHS);
            }

            // 3. Random Shop ID (30% cơ hội là lọc theo Shop)
            if (SHOP_IDS.length > 0 && Math.random() < 0.3) {
                params.shop_id = randomItem(SHOP_IDS);
            }

            // 4. Random Sort (50% cơ hội user sẽ bấm sort)
            if (Math.random() < 0.5) {
                params.sort = randomItem(SORT_OPTIONS);
            }

            // 5. Random Price Range (20% cơ hội user lọc giá)
            if (Math.random() < 0.2) {
                // Tạo giá min ngẫu nhiên từ 0 -> 200k
                const min = Math.floor(Math.random() * 20) * 10000; 
                // Tạo khoảng giá (range) từ 50k -> 500k
                const range = (Math.floor(Math.random() * 10) + 1) * 50000; 
                
                params.price_min = min;
                params.price_max = min + range; // Đảm bảo max luôn > min
            }

            // 6. Random Keywords (10% cơ hội user tìm kiếm từ khóa)
            if (Math.random() < 0.1) {
                const keywords = ["Quần", "Áo", "Giày", "Túi", "Mũ"];
                params.keywords = randomItem(keywords);
            }

            // Chuyển object params thành query string
            const queryString = Object.keys(params)
                .map(key => `${key}=${encodeURIComponent(params[key])}`)
                .join('&');
            
            const res = http.get(`${BASE_URL_PRODUCT}/product/getall?${queryString}`);

            check(res, {
                'List: status 200': (r) => r.status === 200,
                'List: < 800ms': (r) => r.timings.duration < 800,
            });
        });
    } 
    
    // --- SCENARIO B: GET PRODUCT DETAIL - 40% ---
    else if (rand < 0.8) {
        group('API_Get_Product_Detail', function () {
            const key = PRODUCT_KEYS.length > 0 ? randomItem(PRODUCT_KEYS) : 'quan-tap-gym-nam-cao-cap-szone-sq492';
            const res = http.get(`${BASE_URL_PRODUCT}/product/getdetail/${key}`);
            check(res, {
                'Detail: status 200': (r) => r.status === 200,
                'Detail: < 500ms': (r) => r.timings.duration < 500,
            });
        });
    } 
    
    // --- SCENARIO C: CREATE ORDER - 20% ---
    else {
        group('API_Create_Order', function () {

            const items=[];
            for (let i = 0; i < randomInt(1, 5); i++) {
                items.push({
                    sku_id: SKU_IDS.length > 0 ? randomItem(SKU_IDS) : '18bdd27a-1d4b-4d86-90ca-45afd27dbbfc',
                    shop_id: "shop_" + i,
                    quantity: Math.floor(Math.random() * 3) + 1
                });
            }

            const  ACCESS_TOKEN =  generateJWT();
            const payload = JSON.stringify({
                shippingAddress: {
                    fullName: "Load Test User",
                    phone: "0912345678",
                    address: "123 Test Street",
                    district: "District 1",
                    city: "Hanoi",
                    postalCode: "100000"
                },
                paymentMethod: "b2c3d4e5-f6a7-8901-2345-67890abcdef1",
                items: items,
                vouchers: [],
                note: "Test Performance k6 Random"
            });

            const params = { headers: { 'Content-Type': 'application/json', "Authorization": `Bearer ${ACCESS_TOKEN}` } };
            
            // Uncomment dòng dưới để chạy thật
            const res = http.post(`${BASE_URL_ORDER}/orders`, payload, params);
            check(res, {
                'status is 200/201': (r) => r.status === 200 || r.status === 201,
                'latency < 800ms': (r) => r.timings.duration < 800,
            });
        });
    }

    sleep(1);
}

// k6 run --out web-dashboard --out csv=ket_qua_full.csv test.js