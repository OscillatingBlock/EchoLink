import React from 'react'
import DarkVeil from '../../components/Animation/DarkViel'
import logo from '/echolink.png'
import {Link} from 'react-router-dom'

const LandingPage=() =>
{
    return (
        <div className='overflow-x-hidden relative'>

            <div style={{width: '100%', height: '100vh', position: 'relative'}}>
                <DarkVeil />



            </div>
            <div className='absolute top-15 left-[50%] -translate-x-1/2 w-[60%] items-center flex justify-between h-17 rounded-full 
    bg-white/10 backdrop-blur-3xl border border-white/20 shadow-lg px-5 z-50'>
                <div className='h-10 '>
                    <img className='h-full' src={logo} alt="" />
                </div>
                <div className='flex gap-6'>
                    <Link to="/dashboard" className='text-xl'>Home</Link>
                    <Link to="/login" className='text-xl'>Login</Link>
                </div>
            </div>

            <div className='absolute flex flex-col justify-center items-center text-white left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 h-full w-[60%] text-center'>

                <h1 className='text-[3vw] font-semibold'>Build Smart Voice Bots that Answer Real Phone Calls.</h1>
                <div className='mt-8 flex gap-5'>
                    <button className='py-3 px-8 rounded-full bg-white text-black'>Get Started</button>

                    <button className='py-3 px-8 rounded-full bg-[#333] border-amber-50/50 border-2 text-white'>Learn More</button>
                </div>
            </div>


        </div>
    )
}

export default LandingPage