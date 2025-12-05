package gfx11

import "github.com/shinkar94/godesktopdup/interop"

type DeviceChildVtbl struct {
	interop.UnknownVtbl
	GetDevice               uintptr
	GetPrivateData          uintptr
	SetPrivateData          uintptr
	SetPrivateDataInterface uintptr
}

type DeviceContextVtbl struct {
	DeviceChildVtbl
	VSSetConstantBuffers                      uintptr
	PSSetShaderResources                      uintptr
	PSSetShader                               uintptr
	PSSetSamplers                             uintptr
	VSSetShader                               uintptr
	DrawIndexed                               uintptr
	Draw                                      uintptr
	Map                                       uintptr
	Unmap                                     uintptr
	PSSetConstantBuffers                      uintptr
	IASetInputLayout                          uintptr
	IASetVertexBuffers                        uintptr
	IASetIndexBuffer                          uintptr
	DrawIndexedInstanced                      uintptr
	DrawInstanced                             uintptr
	GSSetConstantBuffers                      uintptr
	GSSetShader                               uintptr
	IASetPrimitiveTopology                    uintptr
	VSSetShaderResources                      uintptr
	VSSetSamplers                             uintptr
	Begin                                     uintptr
	End                                       uintptr
	GetData                                   uintptr
	SetPredication                            uintptr
	GSSetShaderResources                      uintptr
	GSSetSamplers                             uintptr
	OMSetRenderTargets                        uintptr
	OMSetRenderTargetsAndUnorderedAccessViews uintptr
	OMSetBlendState                           uintptr
	OMSetDepthStencilState                    uintptr
	SOSetTargets                              uintptr
	DrawAuto                                  uintptr
	DrawIndexedInstancedIndirect              uintptr
	DrawInstancedIndirect                     uintptr
	Dispatch                                  uintptr
	DispatchIndirect                          uintptr
	RSSetState                                uintptr
	RSSetViewports                            uintptr
	RSSetScissorRects                         uintptr
	CopySubresourceRegion                     uintptr
	CopyResource                               uintptr
}

type DeviceVtbl struct {
	interop.UnknownVtbl
	CreateBuffer                         uintptr
	CreateTexture1D                      uintptr
	CreateTexture2D                      uintptr
	CreateTexture3D                      uintptr
	CreateShaderResourceView             uintptr
	CreateUnorderedAccessView            uintptr
	CreateRenderTargetView               uintptr
	CreateDepthStencilView               uintptr
	CreateInputLayout                    uintptr
	CreateVertexShader                   uintptr
	CreateGeometryShader                 uintptr
	CreateGeometryShaderWithStreamOutput uintptr
	CreatePixelShader                    uintptr
	CreateHullShader                      uintptr
	CreateDomainShader                   uintptr
	CreateComputeShader                  uintptr
	CreateClassLinkage                   uintptr
	CreateBlendState                     uintptr
	CreateDepthStencilState              uintptr
	CreateRasterizerState                uintptr
	CreateSamplerState                   uintptr
	CreateQuery                          uintptr
	CreatePredicate                      uintptr
	CreateCounter                        uintptr
	CreateDeferredContext                uintptr
	OpenSharedResource                   uintptr
	CheckFormatSupport                   uintptr
	CheckMultisampleQualityLevels        uintptr
	CheckCounterInfo                     uintptr
	CheckCounter                         uintptr
	CheckFeatureSupport                  uintptr
	GetPrivateData                       uintptr
	SetPrivateData                       uintptr
	SetPrivateDataInterface              uintptr
	GetFeatureLevel                      uintptr
	GetCreationFlags                     uintptr
	GetDeviceRemovedReason               uintptr
	GetImmediateContext                  uintptr
	SetExceptionMode                     uintptr
	GetExceptionMode                     uintptr
}

type Texture2DVtbl struct {
	ResourceVtbl
	GetDesc uintptr
}

type ResourceVtbl struct {
	DeviceChildVtbl
	GetType                       uintptr
	SetEvictionPriority           uintptr
	GetEvictionPriority           uintptr
}

