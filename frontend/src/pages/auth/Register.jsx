import React, {useEffect, useRef} from 'react';
import bot from '/bot.png';
import {useForm} from 'react-hook-form';
import {NavLink, useNavigate} from 'react-router-dom';
import {toast} from 'react-toastify';
import {useDispatch} from "react-redux";
import { registerAndConnect } from '../../features/auth/authSlice';

const Signin=() =>
{
    const dispatch=useDispatch();

    const navigate=useNavigate();
    const isSubmitting=useRef(false);

    const savedForm=localStorage.getItem("signupform");
    const {register, reset, handleSubmit, formState: {errors}, watch}=useForm({
        defaultValues: savedForm? JSON.parse(savedForm):{}
    });
    const values=watch();

    useEffect(() =>
    {
        if (!isSubmitting.current)
        {
            localStorage.setItem("signupform", JSON.stringify(values));
        }
    }, [values]);

    async function SubmitHandler(data)
    {
        try {
          isSubmitting.current=true;
          const res = await dispatch(registerAndConnect(data)).unwrap();
          reset();
          localStorage.removeItem("signupform");
          toast.success(res?.message || "Connected");
          navigate("/dashboard");
        } catch (err) {
          toast.error(err || "Registration failed");
        }
    }

    function backHandler()
    {
        isSubmitting.current=true;
        reset();
        localStorage.removeItem("signupform");
        navigate("/");
    }

    return (
        <div className='min-h-[100vh] p-4 md:p-10 flex flex-col md:flex-row text-white'>
            <div className='linear-gradient relative h-full w-full p-2'>

                {/* Close Button */}
                <div
                    onClick={backHandler}
                    className="absolute cursor-pointer h-6 w-6 md:h-8 md:w-8 flex justify-center items-center right-4 md:right-10 top-4 md:top-10 rounded-full z-80"
                >
                    <i className="ri-close-circle-line bg-gradient-to-r cursor-pointer text-2xl md:text-3xl from-[#c95efa] to-[#803bff] text-transparent bg-clip-text"></i>
                </div>

                <div className='bg-[#111] rounded-2xl relative flex flex-col md:flex-row h-full w-full '>

                    {/* Title */}
                    <div className='absolute top-4 md:top-8 bg-linear-to-r mix-blend-lighten flex justify-center items-center text-xl md:text-3xl z-20 font-semibold w-full'>
                        <div className='px-4 py-2 linear-gradient2'>Join the Future of Digital Calling</div>
                    </div>

                    {/* Left Image */}
                    <div className='h-80 md:h-full w-full pt-20 md:w-2/5 md:trapezium flex justify-center items-center'>
                        <img className='h-full object-contain' src={bot} alt="Bot" />
                    </div>

                    {/* Right Form */}
                    <div className='h-[60%] md:h-full w-full md:w-3/5 p-4 md:p-10 flex flex-col gap-2 justify-center items-center'>
                        <div className='text-xl md:text-2xl mt-4 md:mt-20 font-semibold'>Sign Up</div>

                        <form onSubmit={handleSubmit(SubmitHandler)} className='flex flex-col gap-2 justify-center items-center w-full max-w-[600px]'>

                            <div className='flex gap-4 flex-col h-full w-full mt-2'>

                                {/* First & Last Name */}
                                <div className='w-full flex flex-col md:flex-row gap-2'>
                                    <div className='w-full flex flex-col'>
                                        <input
                                            {...register("first_name", {required: "This is a required field"})}
                                            type="text"
                                            placeholder="First Name"
                                            className="p-2 radius w-full linear-gradient3 text-white border-b-1 placeholder-gray-200"
                                        />
                                        <span className='text-red-400 text-[10px]'>{errors?.first_name?.message}</span>
                                    </div>

                                    <div className='w-full flex flex-col'>
                                        <input
                                            {...register("last_name", {required: "This is a required field"})}
                                            type="text"
                                            placeholder="Last Name"
                                            className="p-2 radius w-full linear-gradient3 text-white border-b-1 placeholder-gray-200"
                                        />
                                        <span className='text-red-400 text-[10px]'>{errors?.last_name?.message}</span>
                                    </div>
                                </div>

                                {/* Email */}
                                <div className='w-full flex flex-col'>
                                    <input
                                        {...register("email", {required: "This is a required field"})}
                                        type="email"
                                        placeholder="Email"
                                        className="p-2 radius linear-gradient3 bg-[#212121] text-white border-b-1 placeholder-gray-200"
                                    />
                                    <span className='text-red-400 text-[10px]'>{errors?.email?.message}</span>
                                </div>

                                {/* Twilio Account SID */}
                                <div className='w-full flex flex-col'>
                                    <input
                                        {...register("account_sid", {required: "This is a required field"})}
                                        type="text"
                                        placeholder="Twilio Account SID"
                                        className="p-2 radius linear-gradient3 bg-[#212121] text-white border-b-1 placeholder-gray-200"
                                    />
                                    <span className='text-red-400 text-[10px]'>{errors?.account_sid?.message}</span>
                                </div>

                                {/* Twilio Auth Token */}
                                <div className='w-full flex flex-col'>
                                    <input
                                        {...register("auth_token", {required: "This is a required field"})}
                                        type="password"
                                        placeholder="Twilio Auth Token"
                                        className="p-2 radius linear-gradient3 bg-[#212121] text-white border-b-1 placeholder-gray-200"
                                    />
                                    <span className='text-red-400 text-[10px]'>{errors?.auth_token?.message}</span>
                                </div>

                                {/* Twilio Phone Number SID */}
                                <div className='w-full flex flex-col'>
                                    <input
                                        {...register("phone_number_sid", {required: "This is a required field"})}
                                        type="text"
                                        placeholder="Twilio Phone Number SID"
                                        className="p-2 radius linear-gradient3 bg-[#212121] text-white border-b-1 placeholder-gray-200"
                                    />
                                    <span className='text-red-400 text-[10px]'>{errors?.phone_number_sid?.message}</span>
                                </div>

                                {/* Submit */}
                                <div className='bg-gradient-to-r text-xl rounded p-[1px] border-purple-400 from-[#c95efa] to-[#803bff]'>
                                    <div className='bg-[#212121] w-full rounded p-2 flex justify-center items-center'>
                                        <button
                                            type="submit"
                                            className="bg-gradient-to-r cursor-pointer w-full from-[#c95efa] to-[#803bff] text-transparent bg-clip-text rounded"
                                        >
                                            Submit
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </form>

                        <div className='text-sm mt-4'>
                            Already have an account?{" "}
                            <NavLink
                                className="bg-gradient-to-r from-[#c95efa] to-[#803bff] bg-clip-text text-transparent"
                                to={"/login"}
                            >
                                Login
                            </NavLink>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Signin;
