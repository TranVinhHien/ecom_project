import API from "./api"

export const handleProductImg = (img: string) => {
return  `${API.base}${API.media.product}?id=${img}`
}
export const handleAvatarImg = (img: string) => {
return  `${API.base}${API.media.avtatar}${img}`
}