import type { ColorPalette } from '@/types'

export const colorPalettes: ColorPalette[] = [
  { id: 'matching-gradient', name: 'Matching Gradient', colors: ['#845ec2', '#2c73d2', '#0081cf', '#0089ba', '#008e9b', '#008f7a'] },
  { id: 'spot', name: 'Spot Palette', colors: ['#845ec2', '#b39cd0', '#fbeaff', '#00c9a7'] },
  { id: 'twisted-spot', name: 'Twisted Spot Palette', colors: ['#845ec2', '#00c9a7', '#c4fcef', '#4d8076'] },
  { id: 'classy', name: 'Classy Palette', colors: ['#845ec2', '#4b4453', '#b0a8b9', '#c34a36', '#ff8066'] },
  { id: 'cube', name: 'Cube Palette', colors: ['#845ec2', '#4e8397', '#d5cabd'] },
  { id: 'switch', name: 'Switch Palette', colors: ['#845ec2', '#f3c5ff', '#00c9a7', '#fefedf'] },
  { id: 'small-switch', name: 'Small Switch Palette', colors: ['#845ec2', '#ddeef5', '#3588a3'] },
  { id: 'skip-gradient', name: 'Skip Gradient', colors: ['#845ec2', '#4ffbdf', '#00c2a8', '#008b74'] },
  { id: 'natural', name: 'Natural Palette', colors: ['#845ec2', '#9b89b3', '#fff6ff', '#fefedf'] },
  { id: 'matching', name: 'Matching Palette', colors: ['#845ec2', '#4b4453', '#b0a8b9', '#00896f', '#00c0a3'] },
  { id: 'squash', name: 'Squash Palette', colors: ['#845ec2', '#926c00', '#008ac4'] },
  { id: 'grey-friends', name: 'Grey Friends', colors: ['#845ec2', '#4b4453', '#b0a8b9'] },
  { id: 'dotting', name: 'Dotting Palette', colors: ['#845ec2', '#b0a8b9', '#c34a36', '#bea6a0'] },
  { id: 'skip-shade', name: 'Skip Shade Gradient', colors: ['#845ec2', '#009efa', '#00d2fc', '#4ffbdf'] },
  { id: 'threedom', name: 'Threedom', colors: ['#845ec2', '#b25b00', '#008d82'] }
]

export const allColors = colorPalettes.flatMap(p => p.colors)
