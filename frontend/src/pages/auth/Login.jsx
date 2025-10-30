import React from 'react';
import Bot1 from '/bot.png';
import {useForm} from 'react-hook-form';
import {NavLink, useNavigate} from 'react-router-dom';
import {toast} from 'react-toastify';

const Login=() =>
{
    const navigate=useNavigate();

    const {register, handleSubmit, formState: {errors}}=useForm();

    const loginHandler=async (data) =>
    {
        try
        {
            const user=JSON.parse(localStorage.getItem("user"));

            if (!user)
            {
                toast.error("No registered user found. Please Sign Up first.");
                return;
            }

            if (data.email===user.email&&data.auth_token===user.auth_token)
            {
                toast.success("Login successful");
                navigate("/dashboard");
            } else
            {
                toast.error("Invalid Email or Auth Token");
            }
        } catch (error)
        {
            console.log(error);
            toast.error("Login failed");
        }
    };

    const backHandler=() =>
    {
        navigate("/");
    };

    return (
        <div className='min-h-[100vh] p-4 md:p-10 flex flex-col md:flex-row text-white'>
            <div className='linear-gradient relative h-full w-full p-2'>

                {/* Close Button */}
                <div
                    onClick={backHandler}
                    className="absolute cursor-pointer h-6 w-6 md:h-8 md:w-8 flex justify-center items-center right-4 md:right-10 top-10 md:top-10 rounded-full z-80"
                >
                    <i className="ri-close-circle-line bg-gradient-to-r text-2xl md:text-3xl from-[#c95efa] to-[#803bff] text-transparent bg-clip-text"></i>
                </div>

                <div className='bg-[#111] pt-10 rounded-2xl relative flex flex-col md:flex-row h-full w-full'>

                    {/* Title */}
                    <div className='absolute left-[50%] translate-x-[-50%] w-[80%] top-4 md:top-8 flex justify-center items-center text-xl md:text-3xl z-20 font-semibold '>
                        <div className='px-4 py-2 linear-gradient2'>
                            Welcome Back to the Future of Agents
                        </div>
                    </div>

                    {/* Left Image Section */}
                    <div className='h-80 md:h-full w-full pt-20 md:w-2/5 md:trapezium flex justify-center items-center'>
                        <div className='h-full w-full flex justify-center'>
                            <img className='h-full object-contain' src={Bot1} alt="Character" />
                        </div>
                    </div>

                    {/* Right Form Section */}
                    <div className='h-[60%] md:h-full w-full md:w-3/5 p-4 md:p-10 flex flex-col gap-2 justify-center items-center'>
                        <div className='text-xl md:text-2xl mt-4 md:mt-20 font-semibold'>Login</div>

                        <form
                            onSubmit={handleSubmit(loginHandler)}
                            className='flex flex-col gap-2 justify-center items-center w-full max-w-[600px]'
                        >
                            <div className='flex flex-col gap-4 w-full mt-4'>

                                {/* Email */}
                                <div className='w-full flex flex-col'>
                                    <input
                                        {...register("email", {required: "Email is required"})}
                                        type="email"
                                        placeholder="Email"
                                        className="p-2 radius linear-gradient3 w-full text-white border-b-1 placeholder-gray-200"
                                    />
                                    <span className='text-red-400 text-xs'>{errors?.email?.message}</span>
                                </div>

                                {/* Auth Token */}
                                <div className='w-full flex flex-col'>
                                    <input
                                        {...register("auth_token", {required: "Auth Token is required"})}
                                        type="password"
                                        placeholder="Twilio Auth Token"
                                        className="p-2 radius linear-gradient3 w-full text-white border-b-1 placeholder-gray-200"
                                    />
                                    <span className='text-red-400 text-xs'>{errors?.auth_token?.message}</span>
                                </div>

                                {/* Submit Button */}
                                <div className='bg-gradient-to-r text-xl rounded p-[1px] border-purple-400 from-[#c95efa] to-[#803bff]'>
                                    <div className='bg-[#212121] w-full rounded p-2 flex justify-center items-center'>
                                        <button
                                            type='submit'
                                            className='bg-gradient-to-r cursor-pointer w-full from-[#c95efa] to-[#803bff] text-transparent bg-clip-text rounded'
                                        >
                                            Login
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </form>

                        {/* Signup Redirect */}
                        <div className='text-sm mt-4'>
                            Don’t have an account?{" "}
                            <NavLink
                                to='/sign-up'
                                className='bg-linear-to-r from-[#c95efa] to-[#803bff] bg-clip-text text-transparent font-semibold'
                            >
                                Sign Up
                            </NavLink>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Login;
