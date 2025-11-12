function splitString(inputString: string) {
    // Loại bỏ dấu ngoặc vuông ở đầu và cuối chuỗi
    const stringWithoutBrackets = inputString.slice(1, -1);

    // Tách chuỗi thành mảng các phần tử
    const itemsArray = stringWithoutBrackets.split(',');

    // Trả về mảng kết quả
    return itemsArray;
}
function getDateNow() {
    const currentDate = new Date();
    return currentDate.getDate().toString().padStart(2, '0') + '/' +
        (currentDate.getMonth() + 1).toString().padStart(2, '0') + '/' +
        currentDate.getFullYear();
}
function dateToString(data: Date) {
    const taskEndDate = new Date(data);

    // Chuyển đổi thành chuỗi "dd/mm/yy"
    return taskEndDate.getDate().toString().padStart(2, '0') + '/' +
        (taskEndDate.getMonth() + 1).toString().padStart(2, '0') + '/' +
        taskEndDate.getFullYear().toString().slice(-2);

}
function trimObjectProperties(obj: any) {
    for (var key in obj) {
        if (typeof obj[key] === 'string') {
            obj[key] = obj[key].trim();
        }
    }
    return obj
}

// write a function to check phone number properties return boolean , check string start with 03...09

function isValidPhoneNumber(phone: string): boolean {
    const phoneRegex = /^(?:\+84|0)(\d{9})$/;
    return phoneRegex.test(phone);
}


// write a function to check email properties return boolean , check string end with @.com

function isValidEmail(email: string): boolean {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
}

export { splitString, getDateNow, dateToString, trimObjectProperties, isValidEmail, isValidPhoneNumber }