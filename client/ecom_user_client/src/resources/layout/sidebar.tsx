"use client"
import React from 'react'
import logo from "../../../public/logo.png"
import Image from 'next/image'

export default function Sidebar() {
    // const [roles, setRoles] = useState<roleE[]>()
    // useEffect(() => {
    //     const loadComponent = async () => {
    //         setRoles(cookies.get<roleE[]>(ROLE_USER));
    //     };
    //     loadComponent();
    // }, [])
    return (
        <div

        >
            <div className='fixed flex border-r-2 flex-col gap-2 h-screen  shadow-sm m-4' style={{ zIndex: 100, minWidth: "20vw" }}>
                {/* <i className="pi pi-align-justify absolute top-0 " style={{ zIndex: 10000, right: "-2.5rem", color: '#708090', fontSize: "2.5rem" }}></i> */}
                <div className='flex justify-center'>
                    <Image src={logo} alt='' width={100} height={100} />
                </div>
                <h2 className='text-center'>Item </h2>
                <ul className='p-2 overflow-y-auto h-full'>
                    {/* {SIDEBAR_MENU.map((item) => <MenuItem key={item.code} item={item} checkPermission={roles} />)} */}
                </ul>
            </div>
        </div>
    )
}
